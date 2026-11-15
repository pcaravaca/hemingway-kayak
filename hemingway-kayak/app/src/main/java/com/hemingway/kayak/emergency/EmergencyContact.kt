package com.hemingway.kayak.emergency

import androidx.room.Entity
import androidx.room.PrimaryKey

/**
 * Contacto de emergencia configurable.
 * Configurable emergency contact.
 */
@Entity(tableName = "emergency_contacts")
data class EmergencyContact(
    @PrimaryKey(autoGenerate = true)
    val id: Int = 0,
    val name: String,
    val phone: String?,
    val email: String?,
    val isEmergencyService: Boolean = false,
    val serviceType: EmergencyServiceType? = null,
    val notificationMethods: Set<NotificationMethod> = setOf(),
    val customMessage: String? = null
)

enum class EmergencyServiceType {
    COAST_GUARD,      // Guardia Costera
    MARINE_RESCUE,    // Rescate Marítimo
    LOCAL_POLICE,     // Policía Local
    AMBULANCE        // Ambulancia
}

enum class NotificationMethod {
    SMS,             // Mensaje de texto
    EMAIL,           // Correo electrónico
    PHONE_CALL,      // Llamada telefónica
    EMERGENCY_API    // API de servicios de emergencia
}
