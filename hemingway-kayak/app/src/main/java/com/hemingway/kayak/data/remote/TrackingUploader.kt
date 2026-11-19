package com.hemingway.kayak.data.remote

import com.hemingway.kayak.location.LocationTracker
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
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

    fun start(tripId: String? = null) {
        tracker.onLocation = { location ->
            scope.launch { tracking.point(location) }
        }
        tracker.startTracking()
        scope.launch { tracking.start(tripId, tracker.lastKnownLocation.value) }
    }

    fun stop() {
        val last = tracker.lastKnownLocation.value
        tracker.onLocation = null
        tracker.stopTracking()
        scope.launch { tracking.stop(last) }
    }
}
