package com.hemingway.kayak.ui.weather

import android.Manifest
import android.content.pm.PackageManager
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.app.ActivityCompat
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import com.google.android.gms.tasks.CancellationTokenSource
import com.hemingway.kayak.R
import com.hemingway.kayak.data.WeatherData
import com.hemingway.kayak.databinding.FragmentWeatherBinding
import dagger.hilt.android.AndroidEntryPoint
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

@AndroidEntryPoint
class WeatherFragment : Fragment() {

    private var _binding: FragmentWeatherBinding? = null
    private val binding get() = _binding!!

    private val viewModel: WeatherViewModel by viewModels()
    private lateinit var fusedLocationClient: FusedLocationProviderClient
    private var currentLatitude = 0.0
    private var currentLongitude = 0.0

    private val locationPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        if (isGranted) {
            getCurrentLocation()
        } else {
            showLocationPermissionError()
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentWeatherBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        fusedLocationClient = LocationServices.getFusedLocationProviderClient(requireActivity())
        setupObservers()
        setupClickListeners()
        checkLocationPermission()
    }

    private fun setupObservers() {
        viewModel.weatherState.observe(viewLifecycleOwner) { state ->
            when (state) {
                is WeatherUiState.Loading -> showLoading()
                is WeatherUiState.Success -> showWeather(state.weatherData, state.isFromCache)
                is WeatherUiState.Error -> showError(state.message, state.cachedData)
            }
        }

        viewModel.isRefreshing.observe(viewLifecycleOwner) { isRefreshing ->
            binding.refreshButton.isEnabled = !isRefreshing
        }
    }

    private fun setupClickListeners() {
        binding.refreshButton.setOnClickListener {
            viewModel.refresh(currentLatitude, currentLongitude)
        }
    }

    private fun checkLocationPermission() {
        when {
            ActivityCompat.checkSelfPermission(
                requireContext(),
                Manifest.permission.ACCESS_FINE_LOCATION
            ) == PackageManager.PERMISSION_GRANTED -> {
                getCurrentLocation()
            }
            else -> {
                locationPermissionLauncher.launch(Manifest.permission.ACCESS_FINE_LOCATION)
            }
        }
    }

    private fun getCurrentLocation() {
        if (ActivityCompat.checkSelfPermission(
                requireContext(),
                Manifest.permission.ACCESS_FINE_LOCATION
            ) != PackageManager.PERMISSION_GRANTED
        ) {
            return
        }

        val cancellationToken = CancellationTokenSource()
        fusedLocationClient.getCurrentLocation(
            Priority.PRIORITY_HIGH_ACCURACY,
            cancellationToken.token
        ).addOnSuccessListener { location ->
            location?.let {
                currentLatitude = it.latitude
                currentLongitude = it.longitude
                binding.locationText.text = String.format(
                    Locale.getDefault(),
                    "%.4f, %.4f",
                    currentLatitude,
                    currentLongitude
                )
                viewModel.loadWeather(currentLatitude, currentLongitude)
            } ?: run {
                binding.locationText.text = getString(R.string.msg_location_required)
                loadDefaultWeather()
            }
        }.addOnFailureListener {
            binding.locationText.text = getString(R.string.msg_location_required)
            loadDefaultWeather()
        }
    }

    private fun loadDefaultWeather() {
        currentLatitude = 40.7128
        currentLongitude = -74.0060
        viewModel.loadWeather(currentLatitude, currentLongitude)
    }

    private fun showLoading() {
        binding.progressBar.isVisible = true
        binding.errorText.isVisible = false
    }

    private fun showWeather(weatherData: WeatherData, isFromCache: Boolean) {
        binding.progressBar.isVisible = false
        binding.errorText.isVisible = false

        binding.temperatureText.text = String.format(Locale.getDefault(), "%.0f°C", weatherData.temperature)
        binding.descriptionText.text = weatherData.description.ifEmpty { "N/A" }
        binding.feelsLikeText.text = getString(R.string.weather_feels_like) + ": ${weatherData.feelsLike}°"
        binding.pressureText.text = "${weatherData.pressure} hPa"
        binding.humidityText.text = "${weatherData.humidity}%"
        binding.windText.text = "${weatherData.windSpeed} km/h ${weatherData.windDirection}"
        binding.cloudsText.text = "${weatherData.cloudCover}%"

        val dateFormat = SimpleDateFormat("MMM dd, HH:mm", Locale.getDefault())
        binding.updateText.text = getString(R.string.weather_update) + ": " + dateFormat.format(weatherData.lastUpdated)

        if (isFromCache) {
            binding.updateText.text = getString(R.string.msg_using_cached_data) + " - " + dateFormat.format(weatherData.lastUpdated)
        }
    }

    private fun showError(message: String, cachedData: WeatherData?) {
        binding.progressBar.isVisible = false
        binding.errorText.isVisible = true
        binding.errorText.text = message

        cachedData?.let { showWeather(it, true) }
    }

    private fun showLocationPermissionError() {
        binding.locationText.text = getString(R.string.location_permission_denied)
        loadDefaultWeather()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
