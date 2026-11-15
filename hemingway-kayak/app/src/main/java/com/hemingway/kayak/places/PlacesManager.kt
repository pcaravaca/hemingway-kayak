package com.hemingway.kayak.places

import android.content.Context
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.osmdroid.util.GeoPoint

class PlacesManager(private val context: Context) {

    private val _currentPlace = MutableLiveData<PlaceInfo?>()
    val currentPlace: LiveData<PlaceInfo?> = _currentPlace

    data class PlaceInfo(
        val name: String,
        val address: String,
        val location: GeoPoint,
        val placeId: String
    )

    suspend fun getCurrentPlace(): PlaceInfo? {
        return withContext(Dispatchers.IO) {
            null
        }
    }

    suspend fun searchNearbyPlaces(query: String): List<PlaceInfo> {
        return withContext(Dispatchers.IO) {
            emptyList()
        }
    }

    fun isPlacesInitialized(): Boolean {
        return false
    }

    companion object {
        fun initialize(context: Context, apiKey: String) {
        }
    }
}
