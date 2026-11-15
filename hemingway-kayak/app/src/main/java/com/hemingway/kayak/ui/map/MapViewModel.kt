package com.hemingway.kayak.ui.map

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.asLiveData
import androidx.lifecycle.viewModelScope
import com.hemingway.kayak.data.FishingSpot
import com.hemingway.kayak.data.FishingSpotDao
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.launch
import org.osmdroid.util.GeoPoint
import javax.inject.Inject

@HiltViewModel
class MapViewModel @Inject constructor(
    private val fishingSpotDao: FishingSpotDao
) : ViewModel() {

    private val _currentLocation = MutableLiveData<GeoPoint?>()
    val currentLocation: LiveData<GeoPoint?> = _currentLocation

    val fishingSpots: LiveData<List<FishingSpot>> = fishingSpotDao.getAllFishingSpots().asLiveData()

    private val _selectedSpot = MutableLiveData<FishingSpot?>()
    val selectedSpot: LiveData<FishingSpot?> = _selectedSpot

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    private val _errorMessage = MutableLiveData<String?>()
    val errorMessage: LiveData<String?> = _errorMessage

    fun setCurrentLocation(geoPoint: GeoPoint) {
        _currentLocation.value = geoPoint
    }

    fun selectSpot(spot: FishingSpot?) {
        _selectedSpot.value = spot
    }

    fun selectFishingSpot(spotId: Int) {
        viewModelScope.launch {
            fishingSpotDao.getFishingSpotById(spotId).collect { spot ->
                _selectedSpot.postValue(spot)
            }
        }
    }

    fun loadFishingSpots() {
        viewModelScope.launch {
            _isLoading.value = true
            try {
                fishingSpotDao.getAllFishingSpots().collect { }
            } catch (e: Exception) {
                _errorMessage.postValue("Error loading spots: ${e.message}")
            } finally {
                _isLoading.value = false
            }
        }
    }

    fun saveFishingSpot(spot: FishingSpot) {
        viewModelScope.launch {
            try {
                fishingSpotDao.insertFishingSpot(spot)
            } catch (e: Exception) {
                _errorMessage.postValue("Error saving spot: ${e.message}")
            }
        }
    }

    fun updateFishingSpot(spot: FishingSpot) {
        viewModelScope.launch {
            try {
                fishingSpotDao.updateFishingSpot(spot)
            } catch (e: Exception) {
                _errorMessage.postValue("Error updating spot: ${e.message}")
            }
        }
    }

    fun deleteFishingSpot(spot: FishingSpot) {
        viewModelScope.launch {
            try {
                fishingSpotDao.deleteFishingSpot(spot)
                if (_selectedSpot.value?.id == spot.id) {
                    _selectedSpot.value = null
                }
            } catch (e: Exception) {
                _errorMessage.postValue("Error deleting spot: ${e.message}")
            }
        }
    }

    fun loadNearbySpots(minLat: Double, maxLat: Double, minLng: Double, maxLng: Double) {
        viewModelScope.launch {
            _isLoading.value = true
            try {
                fishingSpotDao.getNearbyFishingSpots(minLat, maxLat, minLng, maxLng).collect { }
            } catch (e: Exception) {
                _errorMessage.postValue("Error loading nearby spots: ${e.message}")
            } finally {
                _isLoading.value = false
            }
        }
    }

    fun clearError() {
        _errorMessage.value = null
    }
}
