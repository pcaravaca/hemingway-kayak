package com.hemingway.kayak.ui.settings

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.hemingway.kayak.databinding.FragmentSettingsBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class SettingsFragment : Fragment() {

    private var _binding: FragmentSettingsBinding? = null
    private val binding get() = _binding!!

    private val viewModel: SettingsViewModel by viewModels()

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentSettingsBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupObservers()
        setupClickListeners()
        binding.versionText.text = "Version 1.0"
    }

    private fun setupObservers() {
        viewModel.userPreferences.observe(viewLifecycleOwner) { prefs ->
            binding.themeValue.text = if (prefs.darkMode) "Dark" else "Light"
            binding.unitsValue.text = if (prefs.useMetricUnits) "Metric" else "Imperial"
            binding.weatherAlertsSwitch.isChecked = prefs.enableWeatherAlerts
            binding.fishingRecommendationsSwitch.isChecked = prefs.enableFishingRecommendations
            binding.fishingSpotsSwitch.isChecked = prefs.showFishingSpots
            binding.trackRouteSwitch.isChecked = prefs.trackRoute
        }
    }

    private fun setupClickListeners() {
        binding.weatherAlertsSwitch.setOnCheckedChangeListener { _, isChecked ->
            viewModel.updateWeatherAlerts(isChecked)
        }

        binding.fishingRecommendationsSwitch.setOnCheckedChangeListener { _, isChecked ->
            viewModel.updateFishingRecommendations(isChecked)
        }

        binding.fishingSpotsSwitch.setOnCheckedChangeListener { _, isChecked ->
            viewModel.updateShowFishingSpots(isChecked)
        }

        binding.trackRouteSwitch.setOnCheckedChangeListener { _, isChecked ->
            viewModel.updateTrackRoute(isChecked)
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
