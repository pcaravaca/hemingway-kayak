package com.hemingway.kayak.data.remote

import javax.inject.Inject
import javax.inject.Singleton

/** Datos del usuario logueado y sus salidas. */
@Singleton
class UserRepository @Inject constructor(
    private val client: BackendClient
) {
    suspend fun me(): UserDto? =
        runCatching { client.api.me() }.getOrNull()

    suspend fun stats(): StatsDto? =
        runCatching { client.api.myStats() }.getOrNull()

    suspend fun myTrips(): List<TripDto> =
        runCatching { client.api.myTrips() }.getOrDefault(emptyList())
}
