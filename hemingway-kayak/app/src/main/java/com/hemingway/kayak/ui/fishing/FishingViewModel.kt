package com.hemingway.kayak.ui.fishing

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.hemingway.kayak.data.FishingRecommendation
import com.hemingway.kayak.data.Resource
import com.hemingway.kayak.data.WeatherData
import com.hemingway.kayak.data.WeatherRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class FishingViewModel @Inject constructor(
    private val weatherRepository: WeatherRepository
) : ViewModel() {

    private val _fishingState = MutableLiveData<FishingUiState>()
    val fishingState: LiveData<FishingUiState> = _fishingState

    private val _isRefreshing = MutableLiveData<Boolean>()
    val isRefreshing: LiveData<Boolean> = _isRefreshing

    fun loadFishingRecommendation(latitude: Double, longitude: Double) {
        viewModelScope.launch {
            _isRefreshing.value = true
            weatherRepository.getWeatherData(latitude, longitude, forceRefresh = false).collect { resource ->
                when (resource) {
                    is Resource.Loading -> {
                        _fishingState.value = FishingUiState.Loading
                    }
                    is Resource.Success -> {
                        val weatherData = resource.data!!
                        val historicalData = emptyList<WeatherData>()
                        val recommendation = weatherRepository.calculateFishingRecommendation(
                            weatherData,
                            historicalData
                        )
                        _fishingState.value = FishingUiState.Success(
                            weatherData = weatherData,
                            recommendation = recommendation,
                            isFromCache = resource.isFromCache
                        )
                    }
                    is Resource.Error -> {
                        _fishingState.value = FishingUiState.Error(
                            message = resource.message ?: "Unknown error"
                        )
                    }
                }
                _isRefreshing.value = false
            }
        }
    }

    fun refresh(latitude: Double, longitude: Double) {
        viewModelScope.launch {
            _isRefreshing.value = true
            weatherRepository.getWeatherData(latitude, longitude, forceRefresh = true).collect { resource ->
                when (resource) {
                    is Resource.Loading -> {
                        _fishingState.value = FishingUiState.Loading
                    }
                    is Resource.Success -> {
                        val weatherData = resource.data!!
                        val historicalData = emptyList<WeatherData>()
                        val recommendation = weatherRepository.calculateFishingRecommendation(
                            weatherData,
                            historicalData
                        )
                        _fishingState.value = FishingUiState.Success(
                            weatherData = weatherData,
                            recommendation = recommendation,
                            isFromCache = resource.isFromCache
                        )
                    }
                    is Resource.Error -> {
                        _fishingState.value = FishingUiState.Error(
                            message = resource.message ?: "Unknown error"
                        )
                    }
                }
                _isRefreshing.value = false
            }
        }
    }
}

sealed class FishingUiState {
    object Loading : FishingUiState()
    data class Success(
        val weatherData: WeatherData,
        val recommendation: FishingRecommendation,
        val isFromCache: Boolean = false
    ) : FishingUiState()
    data class Error(val message: String) : FishingUiState()
}
