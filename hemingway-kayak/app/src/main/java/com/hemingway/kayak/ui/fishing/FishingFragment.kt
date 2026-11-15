package com.hemingway.kayak.ui.fishing

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
import com.hemingway.kayak.data.FishingRecommendation
import com.hemingway.kayak.data.MoonPhase
import com.hemingway.kayak.data.PressureTrend
import com.hemingway.kayak.data.WeatherData
import com.hemingway.kayak.databinding.FragmentFishingBinding
import dagger.hilt.android.AndroidEntryPoint
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

@AndroidEntryPoint
class FishingFragment : Fragment() {

    private var _binding: FragmentFishingBinding? = null
    private val binding get() = _binding!!

    private val viewModel: FishingViewModel by viewModels()
    private lateinit var fusedLocationClient: FusedLocationProviderClient
    private var currentLatitude = 0.0
    private var currentLongitude = 0.0

    private val locationPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        if (isGranted) {
            getCurrentLocation()
        } else {
            loadDefaultLocation()
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentFishingBinding.inflate(inflater, container, false)
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
        viewModel.fishingState.observe(viewLifecycleOwner) { state ->
            when (state) {
                is FishingUiState.Loading -> showLoading()
                is FishingUiState.Success -> showFishingData(state.weatherData, state.recommendation)
                is FishingUiState.Error -> showError(state.message)
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
                viewModel.loadFishingRecommendation(currentLatitude, currentLongitude)
            } ?: run {
                loadDefaultLocation()
            }
        }.addOnFailureListener {
            loadDefaultLocation()
        }
    }

    private fun loadDefaultLocation() {
        currentLatitude = 40.7128
        currentLongitude = -74.0060
        viewModel.loadFishingRecommendation(currentLatitude, currentLongitude)
    }

    private fun showLoading() {
        binding.progressBar.isVisible = true
        binding.errorText.isVisible = false
    }

    private fun showFishingData(weatherData: WeatherData, recommendation: FishingRecommendation) {
        binding.progressBar.isVisible = false
        binding.errorText.isVisible = false

        binding.scoreText.text = recommendation.favorabilityScore.toString()
        binding.scoreLabelText.text = getScoreLabel(recommendation.favorabilityScore)
        binding.pressureTrendText.text = getPressureTrendLabel(recommendation.pressureTrend)
        binding.moonPhaseText.text = getMoonPhaseLabel(recommendation.moonPhase)
        binding.bestTimeText.text = formatBestTime(recommendation.bestTimeToFish)

        binding.techniquesText.text = recommendation.suggestedTechniques.joinToString("\n") { "• $it" }
        binding.baitsText.text = recommendation.suggestedBaits.joinToString("\n") { "• $it" }
    }

    private fun getScoreLabel(score: Int): String {
        return when {
            score >= 80 -> getString(R.string.fishing_excellent)
            score >= 60 -> getString(R.string.fishing_good)
            score >= 40 -> getString(R.string.fishing_moderate)
            score >= 20 -> getString(R.string.fishing_poor)
            else -> getString(R.string.fishing_very_poor)
        }
    }

    private fun getPressureTrendLabel(trend: PressureTrend): String {
        return when (trend) {
            PressureTrend.RISING -> getString(R.string.pressure_rising)
            PressureTrend.FALLING -> getString(R.string.pressure_falling)
            PressureTrend.STABLE -> getString(R.string.pressure_stable)
        }
    }

    private fun getMoonPhaseLabel(phase: MoonPhase): String {
        return when (phase) {
            MoonPhase.NEW_MOON -> getString(R.string.moon_new)
            MoonPhase.WAXING_CRESCENT -> getString(R.string.moon_waxing_crescent)
            MoonPhase.FIRST_QUARTER -> getString(R.string.moon_first_quarter)
            MoonPhase.WAXING_GIBBOUS -> getString(R.string.moon_waxing_gibbous)
            MoonPhase.FULL_MOON -> getString(R.string.moon_full)
            MoonPhase.WANING_GIBBOUS -> getString(R.string.moon_waning_gibbous)
            MoonPhase.THIRD_QUARTER -> getString(R.string.moon_third_quarter)
            MoonPhase.WANING_CRESCENT -> getString(R.string.moon_waning_crescent)
        }
    }

    private fun formatBestTime(date: Date?): String {
        if (date == null) return "N/A"
        val dateFormat = SimpleDateFormat("MMM dd, HH:mm", Locale.getDefault())
        return dateFormat.format(date)
    }

    private fun showError(message: String) {
        binding.progressBar.isVisible = false
        binding.errorText.isVisible = true
        binding.errorText.text = message
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
