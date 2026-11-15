package com.hemingway.kayak.voice

/**
 * Comandos de voz disponibles en la aplicación.
 * Available voice commands in the application.
 */
enum class VoiceCommand(
    val spanishTriggers: List<String>,
    val englishTriggers: List<String>
) {
    // Navegación / Navigation
    START_NAVIGATION(
        spanishTriggers = listOf("iniciar navegación", "empezar ruta", "navegar"),
        englishTriggers = listOf("start navigation", "begin route", "navigate")
    ),
    STOP_NAVIGATION(
        spanishTriggers = listOf("detener navegación", "parar ruta", "terminar navegación"),
        englishTriggers = listOf("stop navigation", "end route", "finish navigation")
    ),

    // Clima / Weather
    CHECK_WEATHER(
        spanishTriggers = listOf("ver clima", "condiciones actuales", "tiempo"),
        englishTriggers = listOf("check weather", "current conditions", "weather")
    ),
    CHECK_WIND(
        spanishTriggers = listOf("ver viento", "dirección del viento", "velocidad del viento"),
        englishTriggers = listOf("check wind", "wind direction", "wind speed")
    ),

    // Pesca / Fishing
    CHECK_FISHING_CONDITIONS(
        spanishTriggers = listOf("condiciones de pesca", "ver pesca", "analizar pesca"),
        englishTriggers = listOf("fishing conditions", "check fishing", "analyze fishing")
    ),
    MARK_FISHING_SPOT(
        spanishTriggers = listOf("marcar punto", "guardar ubicación", "punto de pesca"),
        englishTriggers = listOf("mark spot", "save location", "fishing spot")
    ),

    // Emergencia / Emergency
    EMERGENCY_SOS(
        spanishTriggers = listOf("emergencia", "ayuda", "sos"),
        englishTriggers = listOf("emergency", "help", "sos")
    ),

    // Control de kayak / Kayak control
    START_TRACKING(
        spanishTriggers = listOf("iniciar seguimiento", "empezar registro", "grabar ruta"),
        englishTriggers = listOf("start tracking", "begin recording", "record route")
    ),
    STOP_TRACKING(
        spanishTriggers = listOf("detener seguimiento", "parar registro", "terminar grabación"),
        englishTriggers = listOf("stop tracking", "end recording", "finish recording")
    ),

    // Información / Information
    CHECK_SPEED(
        spanishTriggers = listOf("velocidad actual", "ver velocidad", "rapidez"),
        englishTriggers = listOf("current speed", "check speed", "velocity")
    ),
    CHECK_DISTANCE(
        spanishTriggers = listOf("distancia recorrida", "ver distancia", "recorrido"),
        englishTriggers = listOf("distance traveled", "check distance", "journey")
    );

    companion object {
        /**
         * Encuentra el comando basado en el texto reconocido.
         * Find command based on recognized text.
         */
        fun fromRecognizedText(text: String, isSpanish: Boolean): VoiceCommand? {
            val normalizedText = text.lowercase().trim()
            return values().find { command ->
                val triggers = if (isSpanish) command.spanishTriggers else command.englishTriggers
                triggers.any { normalizedText.contains(it) }
            }
        }
    }
}
