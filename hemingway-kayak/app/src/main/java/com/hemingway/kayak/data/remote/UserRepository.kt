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

    /** Actualiza los campos editables del perfil (displayName, bio, tracking, avatar). */
    suspend fun update(
        id: String,
        displayName: String?,
        bio: String?,
        trackingEnabled: Boolean,
        avatarUrl: String?
    ): UserDto? = runCatching {
        // El backend responde {success,message}; el perfil actualizado sale de /users/me.
        client.api.updateUser(
            id,
            UpdateUserRequest(displayName, bio, trackingEnabled, avatarUrl)
        )
        client.api.me()
    }.getOrNull()
}
