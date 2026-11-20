package com.hemingway.kayak.data.remote

import android.location.Location
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
 * Puente entre el GPS del equipo y el backend: arranca el tracking, manda cada punto
 * nuevo, crea la salida que lo contiene y la cierra con la distancia acumulada.
 * Sin sesión el backend responde 401 y los llamados se tragan (TrackingRepository ya
 * hace runCatching).
 */
@Singleton
class TrackingUploader @Inject constructor(
    private val tracker: LocationTracker,
    private val tracking: TrackingRepository,
    private val trips: TripRepository
) {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)

    private val _isTracking = MutableStateFlow(false)
    val isTracking: StateFlow<Boolean> = _isTracking.asStateFlow()

    private val _distanceKm = MutableStateFlow(0.0)
    val distanceKm: StateFlow<Double> = _distanceKm.asStateFlow()

    /** La salida creada para esta sesión de tracking, si se pudo crear. */
    private val _tripId = MutableStateFlow<String?>(null)
    val tripId: StateFlow<String?> = _tripId.asStateFlow()

    /** Última ubicación conocida; la usan las pantallas para dibujar. */
    val locations: StateFlow<Location?> get() = tracker.currentLocation

    private var lastLocation: Location? = null
    private var distanceMeters = 0f

    /** true mientras la salida no existe y debe crearse con el primer punto real. */
    private var tripPending = false

    fun start(tripId: String? = null) {
        lastLocation = null
        distanceMeters = 0f
        _distanceKm.value = 0.0
        _tripId.value = tripId
        tripPending = tripId == null

        tracker.onLocation = { location ->
            lastLocation?.let { distanceMeters += it.distanceTo(location) }
            lastLocation = location
            _distanceKm.value = distanceMeters.toDouble() / 1000.0
            if (tripPending) {
                // La salida se crea con el primer fix real: al arrancar, el GPS todavía
                // no tiene posición y la salida quedaba con el start en (0, 0).
                tripPending = false
                scope.launch {
                    val id = trips.createLiveTrip(location.latitude, location.longitude)
                    _tripId.value = id
                    tracking.start(id, location)
                }
            } else {
                scope.launch { tracking.point(location) }
            }
        }
        tracker.startTracking()
        _isTracking.value = true

        // Si ya trajeron un tripId (salida creada desde la web), se abre igual.
        if (!tripPending) {
            val start = tracker.lastKnownLocation.value
            scope.launch { tracking.start(tripId, start) }
        }
    }

    fun stop() {
        val last = tracker.lastKnownLocation.value
        val id = _tripId.value
        tracker.onLocation = null
        tracker.stopTracking()
        scope.launch {
            tracking.stop(last)
            if (id != null) {
                trips.finishTrip(id, last, distanceMeters)
            }
        }
        _isTracking.value = false
    }
}
