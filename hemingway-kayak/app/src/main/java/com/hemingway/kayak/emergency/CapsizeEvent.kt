package com.hemingway.kayak.emergency

import android.location.Location
import androidx.room.Entity
import androidx.room.PrimaryKey
import java.util.*
import com.hemingway.kayak.sensors.EmergencyType

/**
 * Evento de vuelco con información detallada.
 * Detailed capsize event information.
 */
@Entity(tableName = "capsize_events")
data class CapsizeEvent(
    @PrimaryKey(autoGenerate = true)
    val id: Int = 0,
    val timestamp: Long = System.currentTimeMillis(),
    val latitude: Double,
    val longitude: Double,
    val accuracy: Float,
    val waterTemperature: Float? = null,
    val windSpeed: Float? = null,
    val windDirection: Float? = null,
    val batteryLevel: Int? = null,
    val eventType: EmergencyType,
    val notificationsSent: Boolean = false,
    val rescueInitiated: Boolean = false
) {
    fun toLocation(): Location {
        return Location("CapsizeEvent").apply {
            this.latitude = this@CapsizeEvent.latitude
            this.longitude = this@CapsizeEvent.longitude
            this.accuracy = this@CapsizeEvent.accuracy
            this.time = timestamp
        }
    }

    companion object {
        fun fromLocation(location: Location, eventType: EmergencyType): CapsizeEvent {
            return CapsizeEvent(
                latitude = location.latitude,
                longitude = location.longitude,
                accuracy = location.accuracy,
                eventType = eventType
            )
        }
    }
}
