package com.hemingway.kayak.data

import androidx.room.*
import kotlinx.coroutines.flow.Flow
import java.util.Date

/**
 * DAO para operaciones de base de datos con datos meteorológicos.
 * DAO for database operations with weather data.
 */
@Dao
interface WeatherDao {
    /**
     * Insertar nuevos datos meteorológicos.
     * Insert new weather data.
     */
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertWeatherData(weatherData: WeatherData)

    /**
     * Obtener los datos meteorológicos más recientes para una ubicación.
     * Get the most recent weather data for a location.
     */
    @Query("SELECT * FROM weather_data WHERE latitude = :lat AND longitude = :lon ORDER BY timestamp DESC LIMIT 1")
    fun getLatestWeatherData(lat: Double, lon: Double): Flow<WeatherData?>

    /**
     * Obtener el historial de datos meteorológicos para una ubicación.
     * Get weather data history for a location.
     */
    @Query("SELECT * FROM weather_data WHERE latitude = :lat AND longitude = :lon ORDER BY timestamp DESC")
    fun getWeatherHistory(lat: Double, lon: Double): Flow<List<WeatherData>>

    /**
     * Eliminar datos meteorológicos caducados.
     * Delete expired weather data.
     */
    @Query("DELETE FROM weather_data WHERE expirationTime < :currentTime")
    suspend fun deleteExpiredWeatherData(currentTime: Date)

    /**
     * Obtener datos meteorológicos por nombre de ubicación.
     * Get weather data by location name.
     */
    @Query("SELECT * FROM weather_data WHERE locationName LIKE :locationName ORDER BY timestamp DESC LIMIT 1")
    fun getWeatherDataByLocation(locationName: String): Flow<WeatherData?>

    /**
     * Comprobar si existen datos meteorológicos válidos para una ubicación.
     * Check if valid weather data exists for a location.
     */
    @Query("SELECT COUNT(*) FROM weather_data WHERE latitude = :lat AND longitude = :lon AND expirationTime > :currentTime")
    suspend fun hasValidWeatherData(lat: Double, lon: Double, currentTime: Date): Int
}
