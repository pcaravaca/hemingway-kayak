package com.hemingway.kayak.data.remote

import android.location.Location
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Manda el tracking al backend. Si no hay sesión, la API responde 401 y los
 * llamados fallan en silencio: no se quiere tumbar la app por un tracking.
 */
@Singleton
class TrackingRepository @Inject constructor(
    private val client: BackendClient
) {
    suspend fun start(tripId: String?, location: Location?) {
        runCatching {
            client.api.startTracking(
                StartTrackingRequest(
                    tripId = tripId,
                    latitude = location?.latitude ?: 0.0,
                    longitude = location?.longitude ?: 0.0
                )
            )
        }
    }

    suspend fun point(location: Location) {
        runCatching {
            client.api.sendPoint(
                TrackingPointRequest(
                    latitude = location.latitude,
                    longitude = location.longitude,
                    speed = if (location.hasSpeed()) location.speed.toDouble() else null,
                    heading = if (location.hasBearing()) location.bearing.toDouble() else null,
                    altitude = if (location.hasAltitude()) location.altitude else null
                )
            )
        }
    }

    suspend fun stop(location: Location?) {
        runCatching {
            client.api.stopTracking(
                StopTrackingRequest(
                    latitude = location?.latitude ?: 0.0,
                    longitude = location?.longitude ?: 0.0
                )
            )
        }
    }
}
