package com.hemingway.kayak.data.remote

import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Part
import retrofit2.http.Path
import okhttp3.MultipartBody

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

    @GET("api/v1/tournaments/{id}/participants")
    suspend fun participants(@Path("id") id: String): List<ParticipantDto>

    @GET("api/v1/tournaments/{id}/catches")
    suspend fun catches(@Path("id") id: String): List<CatchDto>

    @POST("api/v1/tournaments/{id}/join")
    suspend fun joinTournament(@Path("id") id: String, @Body body: JoinRequest): ApiResponse

    @POST("api/v1/tournaments/{id}/catches")
    suspend fun recordCatch(@Path("id") id: String, @Body body: RecordCatchRequest): ApiResponse

    // --- Usuario y salidas propias ---
    @GET("api/v1/users/me")
    suspend fun me(): UserDto

    @GET("api/v1/users/me/stats")
    suspend fun myStats(): StatsDto

    @GET("api/v1/trips/mine")
    suspend fun myTrips(): List<TripDto>

    @POST("api/v1/trips/")
    suspend fun createTrip(@Body body: TripRequest): TripCreatedDto

    @PUT("api/v1/trips/{id}")
    suspend fun updateTrip(@Path("id") id: String, @Body body: TripUpdateRequest): ApiResponse

    // --- Fotos ---
    @Multipart
    @POST("api/v1/uploads")
    suspend fun uploadPhoto(@Part photo: MultipartBody.Part): UploadDto

    // --- Spots ---
    @GET("api/v1/spots")
    suspend fun spots(): List<SpotDto>

    @POST("api/v1/spots/")
    suspend fun createSpot(@Body body: SpotRequest): ApiResponse

    companion object {
        // Por defecto el emulador; se pisa con API_BASE_URL en local.properties.
        const val DEFAULT_BASE_URL = "http://10.0.2.2:8080/"
    }
}
