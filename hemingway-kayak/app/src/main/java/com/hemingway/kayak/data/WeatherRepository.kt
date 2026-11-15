package com.hemingway.kayak.data

import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.flow.flowOn
import kotlinx.coroutines.withContext
import java.util.Calendar
import java.util.Date
import java.util.concurrent.TimeUnit

/**
 * Repositorio para manejar datos meteorológicos, con soporte para caché offline.
 * Repository to handle weather data, with offline cache support.
 */
class WeatherRepository(private val weatherDao: WeatherDao, private val weatherApiService: WeatherApiService) {
    
    private val TAG = "WeatherRepository"

    /**
     * Obtener datos meteorológicos para una ubicación, con soporte offline.
     * Get weather data for a location, with offline support.
     * 
     * @param latitude Latitud de la ubicación / Location latitude
     * @param longitude Longitud de la ubicación / Location longitude
     * @param forceRefresh Forzar actualización desde API / Force refresh from API
     */
    fun getWeatherData(latitude: Double, longitude: Double, forceRefresh: Boolean = false): Flow<Resource<WeatherData>> = flow {
        emit(Resource.Loading())
        
        // Verificar si tenemos datos válidos en caché
        // Check if we have valid cached data
        val currentTime = Date()
        val hasValidCache = withContext(Dispatchers.IO) {
            weatherDao.hasValidWeatherData(latitude, longitude, currentTime) > 0
        }
        
        // Si tenemos datos válidos en caché y no se fuerza la actualización, usarlos
        // If we have valid cached data and refresh is not forced, use it
        if (hasValidCache && !forceRefresh) {
            val cachedData = weatherDao.getLatestWeatherData(latitude, longitude).firstOrNull()
            if (cachedData != null) {
                Log.d(TAG, "Using cached weather data for location ($latitude, $longitude)")
                emit(Resource.Success(cachedData))
                return@flow
            }
        }
        
        // Si no hay caché válida o se fuerza la actualización, intentar obtener datos de la API
        // If no valid cache or refresh is forced, try to get data from API
        try {
            val apiResponse = weatherApiService.getWeatherData(latitude, longitude)
            
            // Calcular tiempo de expiración (3 horas desde ahora)
            // Calculate expiration time (3 hours from now)
            val calendar = Calendar.getInstance()
            calendar.time = currentTime
            calendar.add(Calendar.HOUR, 3)
            val expirationTime = calendar.time
            
            // Convertir la respuesta API al modelo WeatherData
            // Convert API response to WeatherData model
            val weatherData = apiResponse.toWeatherData(latitude, longitude, expirationTime)
            
            // Guardar en caché
            // Save to cache
            withContext(Dispatchers.IO) {
                weatherDao.insertWeatherData(weatherData)
            }
            
            Log.d(TAG, "Fetched and cached new weather data for location ($latitude, $longitude)")
            emit(Resource.Success(weatherData))
        } catch (e: Exception) {
            Log.e(TAG, "Error fetching weather data: ${e.message}")
            
            // Si falla la API pero tenemos datos en caché (aunque estén caducados), usarlos como respaldo
            // If API fails but we have any cached data (even if expired), use it as fallback
            val cachedData = weatherDao.getLatestWeatherData(latitude, longitude).firstOrNull()
            if (cachedData != null) {
                Log.d(TAG, "Using expired cached data as fallback")
                emit(Resource.Success(cachedData, isFromCache = true))
            } else {
                emit(Resource.Error("No se pudieron cargar los datos meteorológicos: ${e.message}"))
            }
        }
    }.flowOn(Dispatchers.IO)
    
    /**
     * Obtener historial de datos meteorológicos para una ubicación.
     * Get weather data history for a location.
     */
    fun getWeatherHistory(latitude: Double, longitude: Double): Flow<List<WeatherData>> {
        return weatherDao.getWeatherHistory(latitude, longitude)
    }
    
    /**
     * Limpiar datos meteorológicos caducados.
     * Clean up expired weather data.
     */
    suspend fun cleanupExpiredData() {
        withContext(Dispatchers.IO) {
            val currentTime = Date()
            weatherDao.deleteExpiredWeatherData(currentTime)
        }
    }
    
    /**
     * Calcular recomendación de pesca basada en los datos meteorológicos.
     * Calculate fishing recommendation based on weather data.
     */
    fun calculateFishingRecommendation(weatherData: WeatherData, 
                                      historicalData: List<WeatherData>): FishingRecommendation {
        // Determinar tendencia de presión basada en histórico
        // Determine pressure trend based on historical data
        val pressureTrend = when {
            historicalData.size < 2 -> PressureTrend.STABLE
            weatherData.pressure > historicalData[1].pressure + 2 -> PressureTrend.RISING
            weatherData.pressure < historicalData[1].pressure - 2 -> PressureTrend.FALLING
            else -> PressureTrend.STABLE
        }
        
        // Calcular fase lunar (implementación simplificada)
        // Calculate moon phase (simplified implementation)
        val moonPhase = calculateMoonPhase(Date())
        
        // Calcular puntuación de favorabilidad
        // Calculate favorability score
        val favorabilityScore = calculateFavorabilityScore(weatherData, pressureTrend, moonPhase)
        
        // Determinar técnicas y cebos sugeridos
        // Determine suggested techniques and baits
        val suggestedTechniques = getSuggestedTechniques(weatherData, favorabilityScore)
        val suggestedBaits = getSuggestedBaits(weatherData, favorabilityScore)
        
        // Determinar mejor momento para pescar
        // Determine best time to fish
        val bestTimeToFish = calculateBestTimeToFish(weatherData, moonPhase)
        
        return FishingRecommendation(
            favorabilityScore = favorabilityScore,
            pressureTrend = pressureTrend,
            moonPhase = moonPhase,
            bestTimeToFish = bestTimeToFish,
            suggestedTechniques = suggestedTechniques,
            suggestedBaits = suggestedBaits
        )
    }
    
    // Métodos auxiliares para cálculo de recomendaciones
    // Helper methods for recommendation calculations
    
    private fun calculateMoonPhase(date: Date): MoonPhase {
        // Implementación simplificada de cálculo de fase lunar
        // Simplified implementation of moon phase calculation
        val calendar = Calendar.getInstance()
        calendar.time = date
        val dayOfMonth = calendar.get(Calendar.DAY_OF_MONTH)
        
        // Simulación de fase lunar basada en día del mes
        // Moon phase simulation based on day of month
        return when (dayOfMonth % 28) {
            0 -> MoonPhase.NEW_MOON
            in 1..3 -> MoonPhase.WAXING_CRESCENT
            in 4..6 -> MoonPhase.FIRST_QUARTER
            in 7..13 -> MoonPhase.WAXING_GIBBOUS
            14 -> MoonPhase.FULL_MOON
            in 15..20 -> MoonPhase.WANING_GIBBOUS
            in 21..24 -> MoonPhase.THIRD_QUARTER
            else -> MoonPhase.WANING_CRESCENT
        }
    }
    
    private fun calculateFavorabilityScore(
        weatherData: WeatherData,
        pressureTrend: PressureTrend,
        moonPhase: MoonPhase
    ): Int {
        var score = 50 // Puntuación base / Base score
        
        // Factores de ajuste basados en condiciones
        // Adjustment factors based on conditions
        
        // Ajuste por temperatura
        // Temperature adjustment
        when {
            weatherData.temperature in 18f..24f -> score += 15
            weatherData.temperature in 15f..28f -> score += 10
            weatherData.temperature in 10f..30f -> score += 5
            weatherData.temperature < 5f || weatherData.temperature > 35f -> score -= 10
        }
        
        // Ajuste por viento
        // Wind adjustment
        when {
            weatherData.windSpeed < 2f -> score += 10
            weatherData.windSpeed in 2f..5f -> score += 5
            weatherData.windSpeed in 5f..10f -> score -= 5
            weatherData.windSpeed > 10f -> score -= 15
        }
        
        // Ajuste por presión
        // Pressure adjustment
        when (pressureTrend) {
            PressureTrend.STABLE -> score += 10
            PressureTrend.RISING -> score += 5
            PressureTrend.FALLING -> score -= 5
        }
        
        // Ajuste por fase lunar
        // Moon phase adjustment
        when (moonPhase) {
            MoonPhase.FULL_MOON, MoonPhase.NEW_MOON -> score += 15
            MoonPhase.FIRST_QUARTER, MoonPhase.THIRD_QUARTER -> score += 10
            else -> score += 0
        }
        
        // Limitar puntuación entre 0 y 100
        // Limit score between 0 and 100
        return score.coerceIn(0, 100)
    }
    
    private fun getSuggestedTechniques(weatherData: WeatherData, favorabilityScore: Int): List<String> {
        val techniques = mutableListOf<String>()
        
        // Técnicas basadas en condiciones
        // Techniques based on conditions
        when {
            weatherData.windSpeed > 8f -> techniques.add("Casting a sotavento / Casting downwind")
            weatherData.windSpeed < 3f -> techniques.add("Técnicas de superficie / Surface techniques")
        }
        
        when {
            weatherData.cloudCover > 80 -> techniques.add("Señuelos oscuros / Dark colored lures")
            weatherData.cloudCover < 30 -> techniques.add("Señuelos brillantes / Bright colored lures")
        }
        
        // Técnicas basadas en puntuación
        // Techniques based on score
        if (favorabilityScore > 70) {
            techniques.add("Aguas superficiales / Shallow waters")
            techniques.add("Cerca de vegetación / Near vegetation")
        } else if (favorabilityScore < 40) {
            techniques.add("Aguas profundas / Deep waters")
            techniques.add("Pesca lenta / Slow fishing")
        }
        
        return techniques
    }
    
    private fun getSuggestedBaits(weatherData: WeatherData, favorabilityScore: Int): List<String> {
        val baits = mutableListOf<String>()
        
        // Cebos basados en condiciones
        // Baits based on conditions
        when {
            weatherData.temperature < 15f -> baits.add("Jigs / Jigs")
            weatherData.temperature > 25f -> baits.add("Topwater / Topwater")
        }
        
        // Cebos basados en puntuación
        // Baits based on score
        if (favorabilityScore > 70) {
            baits.add("Crankbaits / Crankbaits")
            baits.add("Spinnerbaits / Spinnerbaits")
        } else if (favorabilityScore < 40) {
            baits.add("Soft plastics / Soft plastics")
            baits.add("Carolina rig / Carolina rig")
        }
        
        return baits
    }
    
    private fun calculateBestTimeToFish(weatherData: WeatherData, moonPhase: MoonPhase): Date? {
        val calendar = Calendar.getInstance()
        val currentHour = calendar.get(Calendar.HOUR_OF_DAY)
        
        // Si es después de las 18h, sugerir "mañana temprano"
        // If it's after 6 PM, suggest "early tomorrow"
        if (currentHour >= 18) {
            calendar.add(Calendar.DAY_OF_MONTH, 1)
            calendar.set(Calendar.HOUR_OF_DAY, 6)
            calendar.set(Calendar.MINUTE, 0)
            return calendar.time
        }
        
        // En luna llena o nueva, los mejores momentos son amanecer y atardecer
        // On full moon or new moon, best times are dawn and dusk
        if (moonPhase == MoonPhase.FULL_MOON || moonPhase == MoonPhase.NEW_MOON) {
            return if (currentHour < 12) {
                // Sugerir atardecer
                // Suggest dusk
                calendar.set(Calendar.HOUR_OF_DAY, 18)
                calendar.set(Calendar.MINUTE, 0)
                calendar.time
            } else {
                // Sugerir mañana temprano
                // Suggest early tomorrow
                calendar.add(Calendar.DAY_OF_MONTH, 1)
                calendar.set(Calendar.HOUR_OF_DAY, 6)
                calendar.set(Calendar.MINUTE, 0)
                calendar.time
            }
        }
        
        // Para otras fases, sugerir hora actual + 2 horas
        // For other phases, suggest current time + 2 hours
        calendar.add(Calendar.HOUR, 2)
        return calendar.time
    }
}

/**
 * Clase para encapsular resultados de operaciones de datos.
 * Class to encapsulate data operation results.
 */
sealed class Resource<T>(
    val data: T? = null,
    val message: String? = null,
    val isFromCache: Boolean = false
) {
    class Success<T>(data: T, isFromCache: Boolean = false) : Resource<T>(data, isFromCache = isFromCache)
    class Error<T>(message: String, data: T? = null) : Resource<T>(data, message)
    class Loading<T>(data: T? = null) : Resource<T>(data)
}
