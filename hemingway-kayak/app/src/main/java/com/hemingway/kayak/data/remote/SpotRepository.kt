package com.hemingway.kayak.data.remote

import javax.inject.Inject
import javax.inject.Singleton

/** Spots compartidos: listado y alta. */
@Singleton
class SpotRepository @Inject constructor(
    private val client: BackendClient
) {
    suspend fun all(): List<SpotDto> = try {
        client.api.spots()
    } catch (e: Exception) {
        emptyList()
    }

    suspend fun create(
        name: String,
        description: String?,
        latitude: Double,
        longitude: Double,
        rating: Int
    ): Boolean = try {
        client.api.createSpot(
            SpotRequest(
                name = name,
                description = description,
                latitude = latitude,
                longitude = longitude,
                rating = rating,
                isShared = true
            )
        )
        true
    } catch (e: Exception) {
        false
    }
}
