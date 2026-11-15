package com.hemingway.kayak.data

import com.hemingway.kayak.BuildConfig
import retrofit2.http.GET
import retrofit2.http.Query

interface WeatherApiService {
    @GET("data/2.5/weather")
    suspend fun getWeatherData(
        @Query("lat") lat: Double,
        @Query("lon") lon: Double,
        @Query("appid") appid: String = BuildConfig.WEATHER_API_KEY,
        @Query("units") units: String = "metric",
        @Query("lang") lang: String = "es"
    ): OpenWeatherMapResponse
    
    companion object {
        const val BASE_URL = "https://api.openweathermap.org/"
    }
}
