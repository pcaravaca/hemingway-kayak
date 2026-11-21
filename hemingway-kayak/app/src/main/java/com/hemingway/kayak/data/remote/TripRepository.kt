package com.hemingway.kayak.data.remote

import android.location.Location
import javax.inject.Inject
import javax.inject.Singleton

/** Crea y cierra las salidas propias. */
@Singleton
class TripRepository @Inject constructor(
    private val client: BackendClient
) {
    /**
     * Crea una salida "en vivo" para colgarle el tracking. Devuelve el id o null si no
     * hay sesión o el server falla: el tracking no puede quedarse sin salida, así que
     * quien llama decide si sigue sin tripId.
     */
    suspend fun createLiveTrip(latitude: Double?, longitude: Double?): String? = try {
        val now = java.text.SimpleDateFormat("yyyy-MM-dd HH:mm", java.util.Locale.getDefault())
            .format(java.util.Date())
        client.api.createTrip(
            TripRequest(
                title = "Salida $now",
                description = "Grabada con el tracking en vivo",
                type = "FISHING",
                startLat = latitude,
                startLon = longitude
            )
        ).id
    } catch (e: Exception) {
        null
    }

    /** La salida completa (para el detalle). */
    suspend fun trip(tripId: String): TripDto? =
        runCatching { client.api.trip(tripId) }.getOrNull()

    suspend fun points(tripId: String): List<TripPointDto> =
        runCatching { client.api.tripPoints(tripId) }.getOrDefault(emptyList())

    /** Edita título, descripción, tipo y estado. Devuelve null si salió bien. */
    suspend fun update(
        tripId: String,
        title: String,
        description: String?,
        type: String?,
        status: String?
    ): String? = try {
        client.api.updateTrip(
            tripId,
            TripUpdateRequest(title = title, description = description, type = type, status = status)
        )
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo guardar la salida"
    }

    suspend fun addPoint(tripId: String, latitude: Double, longitude: Double): String? = try {
        client.api.addTripPoint(tripId, TripPointRequest(latitude, longitude))
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo agregar el punto"
    }

    /** Borra la salida (el backend borra los puntos primero). Devuelve null si salió bien. */
    suspend fun delete(tripId: String): String? = try {
        client.api.deleteTrip(tripId)
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo eliminar la salida"
    }

    /** Cierra la salida con la distancia acumulada y la última posición. */
    suspend fun finishTrip(tripId: String, end: Location?, distanceMeters: Float) {
        try {
            client.api.updateTrip(
                tripId,
                TripUpdateRequest(
                    status = "completed",
                    endLat = end?.latitude,
                    endLon = end?.longitude,
                    distance = distanceMeters.toDouble() / 1000.0
                )
            )
        } catch (_: Exception) {
            // El tracking ya se frenó; si no se puede cerrar la salida, no se tumba la app.
        }
    }
}
