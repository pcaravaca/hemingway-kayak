package com.hemingway.kayak.ui.settings

import androidx.lifecycle.LiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.asLiveData
import androidx.lifecycle.viewModelScope
import com.hemingway.kayak.data.UserPreferencesManager
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class SettingsViewModel @Inject constructor(
    private val userPreferencesManager: UserPreferencesManager
) : ViewModel() {

    val userPreferences = userPreferencesManager.userPreferencesFlow.asLiveData()

    fun updateTheme(isDarkMode: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateDarkMode(isDarkMode)
        }
    }

    fun updateUnits(useMetric: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateDistanceUnit(if (useMetric) "metric" else "imperial")
        }
    }

    fun updateWeatherAlerts(enabled: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateWeatherAlerts(enabled)
        }
    }

    fun updateFishingRecommendations(enabled: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateFishingRecommendations(enabled)
        }
    }

    fun updateShowFishingSpots(show: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateShowFishingSpots(show)
        }
    }

    fun updateTrackRoute(track: Boolean) {
        viewModelScope.launch {
            userPreferencesManager.updateTrackRoute(track)
        }
    }

    fun saveEmergencyContact(phone: String, email: String) {
        viewModelScope.launch {
            userPreferencesManager.updateEmergencyContact(phone.trim(), email.trim())
        }
    }
}
