package com.hemingway.kayak.billing

/**
 * Niveles de suscripción disponibles.
 * Available subscription levels.
 */
enum class SubscriptionLevel {
    FREE,           // Básico / Basic
    ANGLER,         // Pescador / Angler
    CAPTAIN,        // Capitán / Captain
    MASTER         // Maestro / Master
}

/**
 * Características incluidas en cada nivel.
 * Features included in each level.
 */
data class SubscriptionFeatures(
    val maxContacts: Int,                  // Contactos máximos / Maximum contacts
    val emergencyServices: Boolean,        // Servicios de emergencia / Emergency services
    val weatherAlerts: Boolean,           // Alertas de clima / Weather alerts
    val offlineMode: Boolean,             // Modo sin conexión / Offline mode
    val customMessages: Boolean,          // Mensajes personalizados / Custom messages
    val prioritySupport: Boolean,         // Soporte prioritario / Priority support
    val adsRemoved: Boolean,              // Sin anuncios / No ads
    val dataExport: Boolean,              // Exportación de datos / Data export
    val updateInterval: Int,              // Intervalo de actualización (ms) / Update interval (ms)
    val fishingSpots: Boolean,            // Puntos de pesca / Fishing spots
    val tideInfo: Boolean,                // Información de mareas / Tide information
    val fishingLog: Boolean,              // Registro de pesca / Fishing log
    val sonarSupport: Boolean,            // Soporte para sonar / Sonar support
    val weatherForecast: Int              // Días de pronóstico / Forecast days
) {
    companion object {
        fun forLevel(level: SubscriptionLevel): SubscriptionFeatures {
            return when (level) {
                SubscriptionLevel.FREE -> SubscriptionFeatures(
                    maxContacts = 3,
                    emergencyServices = false,
                    weatherAlerts = true,
                    offlineMode = true,
                    customMessages = true,
                    prioritySupport = false,
                    adsRemoved = false,
                    dataExport = false,
                    updateInterval = 10000,
                    fishingSpots = true,
                    tideInfo = true,
                    fishingLog = true,
                    sonarSupport = false,
                    weatherForecast = 2
                )
                SubscriptionLevel.ANGLER -> SubscriptionFeatures(
                    maxContacts = 5,
                    emergencyServices = false,
                    weatherAlerts = true,
                    offlineMode = true,
                    customMessages = true,
                    prioritySupport = false,
                    adsRemoved = true,
                    dataExport = false,
                    updateInterval = 5000,
                    fishingSpots = true,
                    tideInfo = true,
                    fishingLog = true,
                    sonarSupport = false,
                    weatherForecast = 3
                )
                SubscriptionLevel.CAPTAIN -> SubscriptionFeatures(
                    maxContacts = 10,
                    emergencyServices = true,
                    weatherAlerts = true,
                    offlineMode = true,
                    customMessages = true,
                    prioritySupport = true,
                    adsRemoved = true,
                    dataExport = true,
                    updateInterval = 3000,
                    fishingSpots = true,
                    tideInfo = true,
                    fishingLog = true,
                    sonarSupport = true,
                    weatherForecast = 7
                )
                SubscriptionLevel.MASTER -> SubscriptionFeatures(
                    maxContacts = Int.MAX_VALUE,
                    emergencyServices = true,
                    weatherAlerts = true,
                    offlineMode = true,
                    customMessages = true,
                    prioritySupport = true,
                    adsRemoved = true,
                    dataExport = true,
                    updateInterval = 1000,
                    fishingSpots = true,
                    tideInfo = true,
                    fishingLog = true,
                    sonarSupport = true,
                    weatherForecast = 14
                )
            }
        }
    }
}
