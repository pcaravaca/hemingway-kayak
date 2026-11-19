package com.hemingway.kayak.data.remote

import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

/**
 * Cliente del backend de Hemingway Kayak (Fiber + PostgreSQL).
 * La URL base sale de BuildConfig.API_BASE_URL.
 */
interface KayakApiService {

    // --- Auth ---
    @POST("api/v1/auth/login")
    suspend fun login(@Body body: AuthRequest): AuthResponse

    @POST("api/v1/auth/register")
    suspend fun register(@Body body: RegisterRequest): AuthResponse

    // --- Tracking ---
    @POST("api/v1/tracking/start")
    suspend fun startTracking(@Body body: StartTrackingRequest): ApiResponse

    @POST("api/v1/tracking/point")
    suspend fun sendPoint(@Body body: TrackingPointRequest): ApiResponse

    @POST("api/v1/tracking/stop")
    suspend fun stopTracking(@Body body: StopTrackingRequest): ApiResponse

    // --- Torneos ---
    @GET("api/v1/tournaments")
    suspend fun tournaments(): List<TournamentDto>

    @GET("api/v1/tournaments/live")
    suspend fun liveTournaments(): List<TournamentDto>

    @GET("api/v1/tournaments/{id}")
    suspend fun tournament(@Path("id") id: String): TournamentDto

    @GET("api/v1/tournaments/{id}/ranking")
    suspend fun ranking(@Path("id") id: String): List<RankingDto>

    @GET("api/v1/tournaments/{id}/catches")
    suspend fun catches(@Path("id") id: String): List<CatchDto>

    @POST("api/v1/tournaments/{id}/join")
    suspend fun joinTournament(@Path("id") id: String, @Body body: JoinRequest): ApiResponse

    @POST("api/v1/tournaments/{id}/catches")
    suspend fun recordCatch(@Path("id") id: String, @Body body: RecordCatchRequest): ApiResponse

    companion object {
        // Por defecto el emulador; se pisa con API_BASE_URL en local.properties.
        const val DEFAULT_BASE_URL = "http://10.0.2.2:8080/"
    }
}
