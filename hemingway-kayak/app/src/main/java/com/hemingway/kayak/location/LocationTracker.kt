package com.hemingway.kayak.location

import android.annotation.SuppressLint
import android.content.Context
import android.location.Location
import android.os.Looper
import com.google.android.gms.location.*
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.tasks.await
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Envía los puntos de ubicación al backend apenas llegan. `LocationTracker` avisa
 * por [onLocation]; si nadie está escuchando, igual actualiza el StateFlow.
 */
@Singleton
class LocationTracker @Inject constructor(
    @ApplicationContext private val context: Context
) {

    private val fusedLocationClient: FusedLocationProviderClient =
        LocationServices.getFusedLocationProviderClient(context)

    private val _currentLocation = MutableStateFlow<Location?>(null)
    val currentLocation: StateFlow<Location?> = _currentLocation

    private val _lastKnownLocation = MutableStateFlow<Location?>(null)
    val lastKnownLocation: StateFlow<Location?> = _lastKnownLocation

    /**
     * Se llama en cada ubicación nueva. Quien arranca el tracking (el backend) se
     * registra acá; si no hay nadie escuchando, se sigue actualizando el StateFlow.
     */
    var onLocation: ((Location) -> Unit)? = null

    private val locationRequest = LocationRequest.create().apply {
        priority = LocationRequest.PRIORITY_HIGH_ACCURACY
        interval = 5000
        fastestInterval = 3000
        smallestDisplacement = 1f
    }

    private val locationCallback = object : LocationCallback() {
        override fun onLocationResult(locationResult: LocationResult) {
            locationResult.lastLocation?.let { location ->
                _currentLocation.value = location
                _lastKnownLocation.value = location
                onLocation?.invoke(location)
            }
        }
    }

    @SuppressLint("MissingPermission")
    fun startTracking() {
        fusedLocationClient.requestLocationUpdates(
            locationRequest,
            locationCallback,
            Looper.getMainLooper()
        )
    }

    fun stopTracking() {
        fusedLocationClient.removeLocationUpdates(locationCallback)
    }

    @SuppressLint("MissingPermission")
    suspend fun getLastLocation(): Location? {
        return try {
            fusedLocationClient.lastLocation.await()
        } catch (e: Exception) {
            null
        }
    }
}
