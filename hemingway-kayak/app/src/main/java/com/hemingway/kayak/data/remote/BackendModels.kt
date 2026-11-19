package com.hemingway.kayak.data.remote

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

// DTOs de la API de Hemingway Kayak. Sólo se mapean los campos que usa la app;
// el backend puede devolver más.

@JsonClass(generateAdapter = true)
data class AuthRequest(
    val email: String,
    val password: String
)

@JsonClass(generateAdapter = true)
data class RegisterRequest(
    val name: String,
    val email: String,
    val username: String,
    val password: String
)

@JsonClass(generateAdapter = true)
data class AuthResponse(
    val token: String,
    val user: UserDto
)

@JsonClass(generateAdapter = true)
data class UserDto(
    val id: String,
    val email: String,
    val username: String?,
    @Json(name = "displayName") val displayName: String?,
    @Json(name = "avatarUrl") val avatarUrl: String?,
    val bio: String?,
    @Json(name = "trackingEnabled") val trackingEnabled: Boolean?
)

@JsonClass(generateAdapter = true)
data class TournamentDto(
    val id: String,
    val name: String,
    val description: String?,
    val status: String?,
    val type: String?,
    @Json(name = "startTime") val startTime: String?,
    @Json(name = "endTime") val endTime: String?,
    @Json(name = "locationName") val locationName: String?,
    @Json(name = "startLat") val startLat: Double?,
    @Json(name = "startLon") val startLon: Double?,
    @Json(name = "organizerUsername") val organizerUsername: String?,
    @Json(name = "participantsCount") val participantsCount: Int?
)

@JsonClass(generateAdapter = true)
data class CatchDto(
    val id: String,
    val species: String?,
    val weight: Double?,
    val username: String?,
    @Json(name = "createdAt") val createdAt: String?
)

@JsonClass(generateAdapter = true)
data class ParticipantDto(
    val id: String,
    @Json(name = "userId") val userId: String?,
    val status: String?,
    val score: Double?,
    @Json(name = "catchesCount") val catchesCount: Long?,
    @Json(name = "displayName") val displayName: String?,
    val username: String?,
    @Json(name = "kayakName") val kayakName: String?
)

@JsonClass(generateAdapter = true)
data class RankingDto(
    val id: String,
    @Json(name = "userId") val userId: String?,
    val rank: Long?,
    val score: Double?,
    @Json(name = "catchesCount") val catchesCount: Long?,
    @Json(name = "displayName") val displayName: String?,
    val username: String?
)

@JsonClass(generateAdapter = true)
data class TrackingPointRequest(
    val latitude: Double,
    val longitude: Double,
    val speed: Double?,
    val heading: Double?,
    val altitude: Double?
)

@JsonClass(generateAdapter = true)
data class StatsDto(
    val sessions: Long?,
    val spots: Long?,
    val tournaments: Long?,
    val trips: Long?,
    @Json(name = "totalDistance") val totalDistance: Double?,
    @Json(name = "totalWaypoints") val totalWaypoints: Long?,
    @Json(name = "liveSessions") val liveSessions: Long?
)

@JsonClass(generateAdapter = true)
data class TripDto(
    val id: String,
    val title: String?,
    val description: String?,
    val type: String?,
    val status: String?,
    @Json(name = "startTime") val startTime: String?,
    @Json(name = "endTime") val endTime: String?,
    val distance: Double?,
    @Json(name = "startLat") val startLat: Double?,
    @Json(name = "startLon") val startLon: Double?
)
