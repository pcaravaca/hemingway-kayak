package com.hemingway.kayak.data

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

/**
 * Modelo de datos para la respuesta de la API de OpenWeatherMap.
 * Data model for OpenWeatherMap API response.
 */
@JsonClass(generateAdapter = true)
data class OpenWeatherMapResponse(
    val coord: Coordinates,
    val weather: List<Weather>,
    val main: Main,
    val wind: Wind,
    val clouds: Clouds,
    val rain: Rain?,
    val snow: Snow?,
    @Json(name = "dt") val timestamp: Long,
    val sys: Sys,
    @Json(name = "id") val cityId: Long,
    @Json(name = "name") val cityName: String,
    @Json(name = "cod") val responseCode: Int
) {
    /**
     * Convertir a modelo de datos de la aplicación.
     * Convert to application data model.
     */
    fun toWeatherData(latitude: Double, longitude: Double, expirationTime: java.util.Date): WeatherData {
        val now = java.util.Date()
        return WeatherData(
            latitude = latitude,
            longitude = longitude,
            locationName = cityName,
            description = weather.firstOrNull()?.description ?: "",
            temperature = main.temp,
            feelsLike = main.feelsLike,
            pressure = main.pressure,
            humidity = main.humidity,
            windSpeed = wind.speed,
            windDirection = wind.deg,
            cloudCover = clouds.all,
            precipitation = (rain?.threeHour ?: 0f) + (snow?.threeHour ?: 0f),
            timestamp = java.util.Date(timestamp * 1000),
            lastUpdated = now,
            expirationTime = expirationTime
        )
    }
}

@JsonClass(generateAdapter = true)
data class Coordinates(
    @Json(name = "lat") val latitude: Double,
    @Json(name = "lon") val longitude: Double
)

@JsonClass(generateAdapter = true)
data class Weather(
    val id: Int,
    val main: String,
    val description: String,
    val icon: String
)

@JsonClass(generateAdapter = true)
data class Main(
    val temp: Float,
    @Json(name = "feels_like") val feelsLike: Float,
    @Json(name = "temp_min") val tempMin: Float,
    @Json(name = "temp_max") val tempMax: Float,
    val pressure: Float,
    val humidity: Int
)

@JsonClass(generateAdapter = true)
data class Wind(
    val speed: Float,
    val deg: Int,
    val gust: Float? = null
)

@JsonClass(generateAdapter = true)
data class Clouds(
    val all: Int
)

@JsonClass(generateAdapter = true)
data class Rain(
    @Json(name = "1h") val oneHour: Float? = null,
    @Json(name = "3h") val threeHour: Float? = null
)

@JsonClass(generateAdapter = true)
data class Snow(
    @Json(name = "1h") val oneHour: Float? = null,
    @Json(name = "3h") val threeHour: Float? = null
)

@JsonClass(generateAdapter = true)
data class Sys(
    val type: Int,
    val id: Int,
    val country: String,
    val sunrise: Long,
    val sunset: Long
)
