package com.hemingway.kayak.ui.settings

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.navigation.fragment.findNavController
import com.hemingway.kayak.R
import com.hemingway.kayak.data.remote.AuthRepository
import com.hemingway.kayak.databinding.FragmentSettingsBinding
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class SettingsFragment : Fragment() {

    private var _binding: FragmentSettingsBinding? = null
    private val binding get() = _binding!!

    private val viewModel: SettingsViewModel by viewModels()

    @Inject
    lateinit var authRepository: AuthRepository

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
            binding.emergencyPhoneInput.setText(prefs.emergencyPhone)
            binding.emergencyEmailInput.setText(prefs.emergencyEmail)
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

        binding.saveEmergencyButton.setOnClickListener {
            val phone = binding.emergencyPhoneInput.text?.toString().orEmpty()
            val email = binding.emergencyEmailInput.text?.toString().orEmpty()
            if (phone.isBlank() && email.isBlank()) {
                binding.savedContactText.text = "Cargá al menos un teléfono o un email"
                binding.savedContactText.isVisible = true
                return@setOnClickListener
            }
            viewModel.saveEmergencyContact(phone, email)
            binding.savedContactText.text = getString(R.string.settings_contact_saved)
            binding.savedContactText.isVisible = true
        }

        binding.logoutButton.setOnClickListener {
            authRepository.logout()
            // Se limpia toda la pila: volver atrás no tiene que devolver a la app.
            findNavController().navigate(
                R.id.loginFragment,
                null,
                androidx.navigation.navOptions {
                    popUpTo(R.id.nav_graph) { inclusive = true }
                }
            )
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
