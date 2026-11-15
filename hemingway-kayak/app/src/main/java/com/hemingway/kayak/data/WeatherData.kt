package com.hemingway.kayak.data

import androidx.room.Entity
import androidx.room.PrimaryKey
import java.util.Date

/**
 * Datos meteorológicos para una ubicación específica.
 * Weather data for a specific location.
 */
@Entity(tableName = "weather_data")
data class WeatherData(
    @PrimaryKey(autoGenerate = true) val id: Int = 0,
    val latitude: Double,
    val longitude: Double,
    val locationName: String,
    val description: String = "",
    val temperature: Float,
    val feelsLike: Float,
    val pressure: Float,
    val humidity: Int,
    val windSpeed: Float,
    val windDirection: Int,
    val cloudCover: Int,
    val precipitation: Float,
    val timestamp: Date,
    val lastUpdated: Date,
    val expirationTime: Date
)

/**
 * Tendencia de presión barométrica.
 * Barometric pressure trend.
 */
enum class PressureTrend {
    RISING,     // Subiendo / Rising
    FALLING,    // Bajando / Falling
    STABLE      // Estable / Stable
}

/**
 * Fase lunar actual.
 * Current moon phase.
 */
enum class MoonPhase {
    NEW_MOON,          // Luna nueva / New moon
    WAXING_CRESCENT,   // Luna creciente / Waxing crescent
    FIRST_QUARTER,     // Cuarto creciente / First quarter
    WAXING_GIBBOUS,    // Gibosa creciente / Waxing gibbous
    FULL_MOON,         // Luna llena / Full moon
    WANING_GIBBOUS,    // Gibosa menguante / Waning gibbous
    THIRD_QUARTER,     // Cuarto menguante / Third quarter
    WANING_CRESCENT    // Luna menguante / Waning crescent
}

/**
 * Recomendación de pesca basada en condiciones meteorológicas y otros factores.
 * Fishing recommendation based on weather conditions and other factors.
 */
data class FishingRecommendation(
    val favorabilityScore: Int,    // Puntuación de 0-100 / Score from 0-100
    val pressureTrend: PressureTrend,
    val moonPhase: MoonPhase,
    val bestTimeToFish: Date?,     // Mejor momento para pescar hoy / Best time to fish today
    val suggestedTechniques: List<String>,  // Técnicas sugeridas / Suggested techniques
    val suggestedBaits: List<String>       // Cebos sugeridos / Suggested baits
)
