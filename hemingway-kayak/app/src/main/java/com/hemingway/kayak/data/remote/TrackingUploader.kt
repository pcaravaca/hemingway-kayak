package com.hemingway.kayak.data.remote

import com.hemingway.kayak.location.LocationTracker
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Puente entre el GPS del equipo y el backend: arranca el tracking, manda cada
 * punto nuevo y lo frena. Sin sesión el backend responde 401 y los llamados se
 * tragan (TrackingRepository ya hace runCatching).
 */
@Singleton
class TrackingUploader @Inject constructor(
    private val tracker: LocationTracker,
    private val tracking: TrackingRepository
) {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)

    private val _isTracking = MutableStateFlow(false)
    val isTracking: StateFlow<Boolean> = _isTracking.asStateFlow()

    /** Última ubicación conocida; la usan las pantallas para dibujar. */
    val locations: StateFlow<android.location.Location?> get() = tracker.currentLocation

    fun start(tripId: String? = null) {
        tracker.onLocation = { location ->
            scope.launch { tracking.point(location) }
        }
        tracker.startTracking()
        scope.launch { tracking.start(tripId, tracker.lastKnownLocation.value) }
        _isTracking.value = true
    }

    fun stop() {
        val last = tracker.lastKnownLocation.value
        tracker.onLocation = null
        tracker.stopTracking()
        scope.launch { tracking.stop(last) }
        _isTracking.value = false
    }
}
