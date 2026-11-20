package com.hemingway.kayak.data

import android.content.Context
import dagger.hilt.android.qualifiers.ApplicationContext
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.floatPreferencesKey
import androidx.datastore.preferences.core.intPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Extension para crear DataStore para preferencias de usuario.
 * Extension to create DataStore for user preferences.
 */
private val Context.userPreferencesDataStore: DataStore<Preferences> by preferencesDataStore(
    name = "user_preferences"
)

/**
 * Clase para gestionar las preferencias del usuario.
 * Class to manage user preferences.
 */
@Singleton
class UserPreferencesManager @Inject constructor(
    @ApplicationContext
    private val context: Context
) {
    private val dataStore = context.userPreferencesDataStore

    // Claves para las preferencias / Preference keys
    companion object {
        // Preferencias generales / General preferences
        val TEMPERATURE_UNIT = stringPreferencesKey("temperature_unit")
        val DISTANCE_UNIT = stringPreferencesKey("distance_unit")
        val LANGUAGE = stringPreferencesKey("language")
        val DARK_MODE = booleanPreferencesKey("dark_mode")
        
        // Preferencias de pesca / Fishing preferences
        val FISH_SPECIES = stringPreferencesKey("fish_species")
        val FISHING_TECHNIQUE = stringPreferencesKey("fishing_technique")
        val PREFERRED_BAIT_TYPES = stringPreferencesKey("preferred_bait_types")
        val PREFERRED_FISHING_TIME = stringPreferencesKey("preferred_fishing_time")
        
        // Preferencias de notificaciones / Notification preferences
        val ENABLE_WEATHER_ALERTS = booleanPreferencesKey("enable_weather_alerts")
        val ENABLE_FISHING_RECOMMENDATIONS = booleanPreferencesKey("enable_fishing_recommendations")
        val NOTIFICATION_FREQUENCY = intPreferencesKey("notification_frequency")
        val MINIMUM_FAVORABILITY_SCORE = intPreferencesKey("minimum_favorability_score")
        
        // Preferencias de mapa y navegación / Map and navigation preferences
        val DEFAULT_MAP_ZOOM = floatPreferencesKey("default_map_zoom")
        val SHOW_FISHING_SPOTS = booleanPreferencesKey("show_fishing_spots")
        val TRACK_ROUTE = booleanPreferencesKey("track_route")
        val ENABLE_OFFLINE_MAPS = booleanPreferencesKey("enable_offline_maps")

        // Contacto de emergencia / Emergency contact
        val EMERGENCY_PHONE = stringPreferencesKey("emergency_phone")
        val EMERGENCY_EMAIL = stringPreferencesKey("emergency_email")
    }

    /**
     * Modelo de datos para las preferencias de usuario.
     * Data model for user preferences.
     */
    data class UserPreferences(
        // Preferencias generales / General preferences
        val temperatureUnit: String = "celsius", // celsius, fahrenheit
        val distanceUnit: String = "metric", // metric, imperial
        val useMetricUnits: Boolean = true,
        val language: String = "es", // es, en
        val darkMode: Boolean = true,
        
        // Preferencias de pesca / Fishing preferences
        val fishSpecies: String = "bass",
        val fishingTechnique: String = "casting",
        val preferredBaitTypes: String = "artificial,live",
        val preferredFishingTime: String = "morning",
        
        // Preferencias de notificaciones / Notification preferences
        val enableWeatherAlerts: Boolean = true,
        val enableFishingRecommendations: Boolean = true,
        val notificationFrequency: Int = 6, // Horas / Hours
        val minimumFavorabilityScore: Int = 70,
        
        // Preferencias de mapa y navegación / Map and navigation preferences
        val defaultMapZoom: Float = 14f,
        val showFishingSpots: Boolean = true,
        val trackRoute: Boolean = true,
        val enableOfflineMaps: Boolean = true,

        // Contacto de emergencia / Emergency contact
        val emergencyPhone: String = "",
        val emergencyEmail: String = ""
    )

    /**
     * Obtener las preferencias del usuario como flujo observable.
     * Get user preferences as an observable flow.
     */
    val userPreferencesFlow: Flow<UserPreferences> = dataStore.data.map { preferences ->
        UserPreferences(
            // Preferencias generales / General preferences
            temperatureUnit = preferences[TEMPERATURE_UNIT] ?: "celsius",
            distanceUnit = preferences[DISTANCE_UNIT] ?: "metric",
            useMetricUnits = preferences[DISTANCE_UNIT] ?: "metric" == "metric",
            language = preferences[LANGUAGE] ?: "es",
            darkMode = preferences[DARK_MODE] ?: true,
            
            // Preferencias de pesca / Fishing preferences
            fishSpecies = preferences[FISH_SPECIES] ?: "bass",
            fishingTechnique = preferences[FISHING_TECHNIQUE] ?: "casting",
            preferredBaitTypes = preferences[PREFERRED_BAIT_TYPES] ?: "artificial,live",
            preferredFishingTime = preferences[PREFERRED_FISHING_TIME] ?: "morning",
            
            // Preferencias de notificaciones / Notification preferences
            enableWeatherAlerts = preferences[ENABLE_WEATHER_ALERTS] ?: true,
            enableFishingRecommendations = preferences[ENABLE_FISHING_RECOMMENDATIONS] ?: true,
            notificationFrequency = preferences[NOTIFICATION_FREQUENCY] ?: 6,
            minimumFavorabilityScore = preferences[MINIMUM_FAVORABILITY_SCORE] ?: 70,
            
            // Preferencias de mapa y navegación / Map and navigation preferences
            defaultMapZoom = preferences[DEFAULT_MAP_ZOOM] ?: 14f,
            showFishingSpots = preferences[SHOW_FISHING_SPOTS] ?: true,
            trackRoute = preferences[TRACK_ROUTE] ?: true,
            enableOfflineMaps = preferences[ENABLE_OFFLINE_MAPS] ?: true,

            emergencyPhone = preferences[EMERGENCY_PHONE] ?: "",
            emergencyEmail = preferences[EMERGENCY_EMAIL] ?: ""
        )
    }
    
    /**
     * Actualizar unidad de temperatura.
     * Update temperature unit.
     */
    suspend fun updateTemperatureUnit(unit: String) {
        dataStore.edit { preferences ->
            preferences[TEMPERATURE_UNIT] = unit
        }
    }
    
    /**
     * Actualizar unidad de distancia.
     * Update distance unit.
     */
    suspend fun updateDistanceUnit(unit: String) {
        dataStore.edit { preferences ->
            preferences[DISTANCE_UNIT] = unit
        }
    }
    
    /**
     * Actualizar idioma.
     * Update language.
     */
    suspend fun updateLanguage(language: String) {
        dataStore.edit { preferences ->
            preferences[LANGUAGE] = language
        }
    }
    
    /**
     * Actualizar modo oscuro.
     * Update dark mode.
     */
    suspend fun updateDarkMode(enabled: Boolean) {
        dataStore.edit { preferences ->
            preferences[DARK_MODE] = enabled
        }
    }
    
    /**
     * Actualizar especies de peces.
     * Update fish species.
     */
    suspend fun updateFishSpecies(species: String) {
        dataStore.edit { preferences ->
            preferences[FISH_SPECIES] = species
        }
    }
    
    /**
     * Actualizar técnica de pesca.
     * Update fishing technique.
     */
    suspend fun updateFishingTechnique(technique: String) {
        dataStore.edit { preferences ->
            preferences[FISHING_TECHNIQUE] = technique
        }
    }
    
    /**
     * Actualizar tipos de cebos preferidos.
     * Update preferred bait types.
     */
    suspend fun updatePreferredBaitTypes(baitTypes: String) {
        dataStore.edit { preferences ->
            preferences[PREFERRED_BAIT_TYPES] = baitTypes
        }
    }
    
    /**
     * Actualizar hora de pesca preferida.
     * Update preferred fishing time.
     */
    suspend fun updatePreferredFishingTime(time: String) {
        dataStore.edit { preferences ->
            preferences[PREFERRED_FISHING_TIME] = time
        }
    }
    
    /**
     * Actualizar alertas meteorológicas.
     * Update weather alerts.
     */
    suspend fun updateWeatherAlerts(enabled: Boolean) {
        dataStore.edit { preferences ->
            preferences[ENABLE_WEATHER_ALERTS] = enabled
        }
    }
    
    /**
     * Actualizar recomendaciones de pesca.
     * Update fishing recommendations.
     */
    suspend fun updateFishingRecommendations(enabled: Boolean) {
        dataStore.edit { preferences ->
            preferences[ENABLE_FISHING_RECOMMENDATIONS] = enabled
        }
    }
    
    /**
     * Actualizar frecuencia de notificaciones.
     * Update notification frequency.
     */
    suspend fun updateNotificationFrequency(hours: Int) {
        dataStore.edit { preferences ->
            preferences[NOTIFICATION_FREQUENCY] = hours
        }
    }
    
    /**
     * Actualizar puntuación mínima de favorabilidad.
     * Update minimum favorability score.
     */
    suspend fun updateMinimumFavorabilityScore(score: Int) {
        dataStore.edit { preferences ->
            preferences[MINIMUM_FAVORABILITY_SCORE] = score
        }
    }
    
    /**
     * Actualizar zoom predeterminado del mapa.
     * Update default map zoom.
     */
    suspend fun updateDefaultMapZoom(zoom: Float) {
        dataStore.edit { preferences ->
            preferences[DEFAULT_MAP_ZOOM] = zoom
        }
    }
    
    /**
     * Actualizar visualización de lugares de pesca.
     * Update fishing spots display.
     */
    suspend fun updateShowFishingSpots(show: Boolean) {
        dataStore.edit { preferences ->
            preferences[SHOW_FISHING_SPOTS] = show
        }
    }
    
    /**
     * Actualizar seguimiento de ruta.
     * Update route tracking.
     */
    suspend fun updateTrackRoute(track: Boolean) {
        dataStore.edit { preferences ->
            preferences[TRACK_ROUTE] = track
        }
    }
    
    /**
     * Actualizar mapas sin conexión.
     * Update offline maps.
     */
    suspend fun updateEnableOfflineMaps(enable: Boolean) {
        dataStore.edit { preferences ->
            preferences[ENABLE_OFFLINE_MAPS] = enable
        }
    }

    /**
     * Guardar el contacto de emergencia (telefono y/o email).
     * Save the emergency contact (phone and/or email).
     */
    suspend fun updateEmergencyContact(phone: String, email: String) {
        dataStore.edit { preferences ->
            preferences[EMERGENCY_PHONE] = phone
            preferences[EMERGENCY_EMAIL] = email
        }
    }
}
