package com.hemingway.kayak.ui.weather

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.asLiveData
import androidx.lifecycle.viewModelScope
import com.hemingway.kayak.data.Resource
import com.hemingway.kayak.data.UserPreferencesManager
import com.hemingway.kayak.data.WeatherData
import com.hemingway.kayak.data.WeatherRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class WeatherViewModel @Inject constructor(
    private val weatherRepository: WeatherRepository,
    private val userPreferencesManager: UserPreferencesManager
) : ViewModel() {

    private val _weatherState = MutableLiveData<WeatherUiState>()
    val weatherState: LiveData<WeatherUiState> = _weatherState

    private val _isRefreshing = MutableLiveData<Boolean>()
    val isRefreshing: LiveData<Boolean> = _isRefreshing

    val useMetricUnits: LiveData<Boolean> = userPreferencesManager.userPreferencesFlow
        .map { it.useMetricUnits }
        .asLiveData()

    fun loadWeather(latitude: Double, longitude: Double, forceRefresh: Boolean = false) {
        viewModelScope.launch {
            _isRefreshing.value = true
            weatherRepository.getWeatherData(latitude, longitude, forceRefresh).collect { resource ->
                when (resource) {
                    is Resource.Loading -> {
                        _weatherState.value = WeatherUiState.Loading
                    }
                    is Resource.Success -> {
                        _weatherState.value = WeatherUiState.Success(
                            weatherData = resource.data!!,
                            isFromCache = resource.isFromCache
                        )
                    }
                    is Resource.Error -> {
                        _weatherState.value = WeatherUiState.Error(
                            message = resource.message ?: "Unknown error",
                            cachedData = resource.data
                        )
                    }
                }
                if (forceRefresh) {
                    _isRefreshing.value = false
                }
            }
        }
        _isRefreshing.value = false
    }

    fun refresh(latitude: Double, longitude: Double) {
        loadWeather(latitude, longitude, forceRefresh = true)
    }
}

sealed class WeatherUiState {
    object Loading : WeatherUiState()
    data class Success(val weatherData: WeatherData, val isFromCache: Boolean = false) : WeatherUiState()
    data class Error(val message: String, val cachedData: WeatherData? = null) : WeatherUiState()
}
