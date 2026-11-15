package com.hemingway.kayak.emergency

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import android.telephony.SmsManager
import android.widget.Toast
import androidx.core.content.ContextCompat
import com.hemingway.kayak.sensors.EmergencyType
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData

class EmergencyNotifier(private val context: Context) {

    private val _emergencyState = MutableLiveData<EmergencyState>()
    val emergencyState: LiveData<EmergencyState> = _emergencyState

    enum class EmergencyState {
        IDLE,
        TRIGGERED,
        SENDING,
        SENT,
        FAILED
    }

    data class EmergencyInfo(
        val latitude: Double,
        val longitude: Double,
        val timestamp: Long,
        val type: String,
        val accuracy: Float = 0f,
        val driftDistance: Float = 0f,
        val driftDirection: Float = 0f,
        val driftSpeed: Float = 0f,
        val waterTemperature: Float? = null,
        val windSpeed: Float? = null,
        val windDirection: Float? = null,
        val batteryLevel: Int? = null,
        val customMessage: String? = null
    ) {
        val compassDirection: String
            get() = when {
                driftDirection < 22.5 || driftDirection >= 337.5 -> "N"
                driftDirection < 67.5 -> "NE"
                driftDirection < 112.5 -> "E"
                driftDirection < 157.5 -> "SE"
                driftDirection < 202.5 -> "S"
                driftDirection < 247.5 -> "SW"
                driftDirection < 292.5 -> "W"
                else -> "NW"
            }
    }

    fun sendEmergencyAlert(type: EmergencyType) {
        val info = EmergencyInfo(0.0, 0.0, System.currentTimeMillis(), type.name)
        sendEmergencyAlert(info)
    }

    fun sendEmergencyAlert(info: EmergencyInfo) {
        _emergencyState.value = EmergencyState.TRIGGERED
    }

    fun sendSmsAlert(info: EmergencyInfo, phoneNumbers: List<String>) {
        if (ContextCompat.checkSelfPermission(context, Manifest.permission.SEND_SMS) 
            != PackageManager.PERMISSION_GRANTED) {
            _emergencyState.value = EmergencyState.FAILED
            return
        }

        _emergencyState.value = EmergencyState.SENDING
        
        try {
            val message = buildEmergencyMessage(info)
            val smsManager = context.getSystemService(SmsManager::class.java)
            
            phoneNumbers.forEach { phoneNumber ->
                val parts = smsManager.divideMessage(message)
                smsManager.sendMultipartTextMessage(
                    phoneNumber.trim(),
                    null,
                    parts,
                    null,
                    null
                )
            }
            
            _emergencyState.value = EmergencyState.SENT
        } catch (e: Exception) {
            _emergencyState.value = EmergencyState.FAILED
        }
    }

    fun sendEmailAlert(info: EmergencyInfo, email: String) {
        _emergencyState.value = EmergencyState.SENDING
        
        val intent = android.content.Intent(android.content.Intent.ACTION_SENDTO).apply {
            data = android.net.Uri.parse("mailto:")
            putExtra(android.content.Intent.EXTRA_EMAIL, arrayOf(email))
            putExtra(android.content.Intent.EXTRA_SUBJECT, "EMERGENCY: Kayak Alert")
            putExtra(android.content.Intent.EXTRA_TEXT, buildEmergencyMessage(info))
            addFlags(android.content.Intent.FLAG_ACTIVITY_NEW_TASK)
        }
        
        try {
            context.startActivity(intent)
            _emergencyState.value = EmergencyState.SENT
        } catch (e: Exception) {
            _emergencyState.value = EmergencyState.FAILED
        }
    }

    fun callEmergencyServices(phoneNumber: String) {
        val intent = android.content.Intent(android.content.Intent.ACTION_DIAL).apply {
            data = android.net.Uri.parse("tel:$phoneNumber")
            addFlags(android.content.Intent.FLAG_ACTIVITY_NEW_TASK)
        }
        context.startActivity(intent)
    }

    fun call911() {
        callEmergencyServices("911")
    }

    private fun buildEmergencyMessage(info: EmergencyInfo): String {
        val timestamp = java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(java.util.Date(info.timestamp))
        val mapsUrl = "https://maps.google.com/?q=${info.latitude},${info.longitude}"
        
        return """
EMERGENCY ALERT / ALERTA DE EMERGENCIA

Type / Tipo: ${info.type}
Time / Hora: $timestamp

Location: $mapsUrl

Drift:
  - Distance / Distancia: ${info.driftDistance.toInt()}m
  - Direction / Dirección: ${info.compassDirection}
  - Speed / Velocidad: ${info.driftSpeed} m/s

${info.customMessage ?: ""}
        """.trimIndent()
    }

    fun buildCapsizeMessage(
        latitude: Double,
        longitude: Double,
        accuracy: Float,
        driftDistance: Float,
        driftDirection: Float,
        driftSpeed: Float,
        waterTemperature: Float?,
        windSpeed: Float?,
        windDirection: Float?,
        batteryLevel: Int?,
        customMessage: String?
    ): String {
        val timestamp = java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(java.util.Date())
        val mapsUrl = "https://maps.google.com/?q=$latitude,$longitude"
        
        val compassDirection = when {
            driftDirection < 22.5 || driftDirection >= 337.5 -> "N"
            driftDirection < 67.5 -> "NE"
            driftDirection < 112.5 -> "E"
            driftDirection < 157.5 -> "SE"
            driftDirection < 202.5 -> "S"
            driftDirection < 247.5 -> "SW"
            driftDirection < 292.5 -> "W"
            else -> "NW"
        }
        
        return """
EMERGENCY ALERT / ALERTA DE EMERGENCIA

$customMessage

Type / Tipo: CAPSIZE
Time / Hora: $timestamp

Location: $mapsUrl

Drift:
  - Distance: ${driftDistance.toInt()}m
  - Direction: $compassDirection
  - Speed: $driftSpeed m/s

Conditions:
  - Water temp: ${waterTemperature ?: "N/A"}°C
  - Wind: ${windSpeed ?: "N/A"} m/s
  - Battery: ${batteryLevel ?: "N/A"}%
        """.trimIndent()
    }

    fun setEmergencyContacts(contacts: List<EmergencyContact>) {
    }

    fun triggerCapsizeAlert(capsizeEvent: CapsizeEvent) {
        val info = EmergencyInfo(
            latitude = capsizeEvent.latitude,
            longitude = capsizeEvent.longitude,
            timestamp = capsizeEvent.timestamp,
            type = "CAPSIZE"
        )
        sendEmergencyAlert(info)
    }

    companion object {
        fun create(context: Context): EmergencyNotifier {
            return EmergencyNotifier(context)
        }
    }
}
