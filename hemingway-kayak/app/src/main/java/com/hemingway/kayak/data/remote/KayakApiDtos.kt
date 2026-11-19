package com.hemingway.kayak.data.remote

import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ApiResponse(val success: Boolean?, val message: String?, val error: String?)

@JsonClass(generateAdapter = true)
data class StartTrackingRequest(
    val tripId: String? = null,
    val latitude: Double,
    val longitude: Double
)

@JsonClass(generateAdapter = true)
data class StopTrackingRequest(
    val latitude: Double,
    val longitude: Double
)

@JsonClass(generateAdapter = true)
data class JoinRequest(
    val kayakId: String? = null
)

@JsonClass(generateAdapter = true)
data class RecordCatchRequest(
    val participantId: String,
    val species: String,
    val weight: Double,
    val latitude: Double?,
    val longitude: Double?,
    val photoUrl: String?
)
