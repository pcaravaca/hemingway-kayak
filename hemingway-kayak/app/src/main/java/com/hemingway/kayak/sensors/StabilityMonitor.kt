package com.hemingway.kayak.sensors

import android.content.Context
import android.hardware.Sensor
import android.hardware.SensorEvent
import android.hardware.SensorEventListener
import android.hardware.SensorManager
import android.media.AudioAttributes
import android.media.AudioManager
import android.media.ToneGenerator
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import kotlin.math.abs
import kotlin.math.atan2
import kotlin.math.sqrt
import com.hemingway.kayak.emergency.EmergencyNotifier
import com.hemingway.kayak.sensors.EmergencyType

/**
 * Monitor de estabilidad del kayak usando sensores del dispositivo.
 * Kayak stability monitor using device sensors.
 */
class StabilityMonitor(private val context: Context) : SensorEventListener {

    private val sensorManager = context.getSystemService(Context.SENSOR_SERVICE) as SensorManager
    private val accelerometer = sensorManager.getDefaultSensor(Sensor.TYPE_ACCELEROMETER)
    private val gyroscope = sensorManager.getDefaultSensor(Sensor.TYPE_GYROSCOPE)

    private val _stability = MutableLiveData<StabilityStatus>()
    val stability: LiveData<StabilityStatus> = _stability

    private val toneGenerator = ToneGenerator(AudioManager.STREAM_ALARM, 100)

    // Umbrales de inclinación / Tilt thresholds
    companion object {
        private const val WARNING_ANGLE = 20f  // Grados / Degrees
        private const val DANGER_ANGLE = 35f   // Grados / Degrees
        private const val CAPSIZED_ANGLE = 90f // Grados / Degrees
        private const val SUBMERGED_TIMEOUT = 30_000L // 30 segundos / 30 seconds
        private const val SAMPLING_PERIOD = 100_000 // Microsegundos / Microseconds
    }

    /**
     * Estado de estabilidad del kayak.
     * Kayak stability status.
     */
    enum class StabilityStatus {
        STABLE,        // Estable / Stable
        WARNING,       // Precaución / Warning
        DANGER,        // Peligro / Danger
        CAPSIZED,      // Volcado / Capsized
        SUBMERGED      // Hundido / Submerged
    }

    private var capsizedTime: Long = 0
    private var emergencyNotifier: EmergencyNotifier? = null

    fun setEmergencyNotifier(notifier: EmergencyNotifier) {
        emergencyNotifier = notifier
    }

    /**
     * Iniciar monitoreo.
     * Start monitoring.
     */
    fun startMonitoring() {
        sensorManager.registerListener(
            this,
            accelerometer,
            SAMPLING_PERIOD
        )
        sensorManager.registerListener(
            this,
            gyroscope,
            SAMPLING_PERIOD
        )
    }

    /**
     * Detener monitoreo.
     * Stop monitoring.
     */
    fun stopMonitoring() {
        sensorManager.unregisterListener(this)
    }

    override fun onSensorChanged(event: SensorEvent) {
        when (event.sensor.type) {
            Sensor.TYPE_ACCELEROMETER -> processAccelerometerData(event)
            Sensor.TYPE_GYROSCOPE -> processGyroscopeData(event)
        }
    }

    private fun processAccelerometerData(event: SensorEvent) {
        val x = event.values[0]
        val y = event.values[1]
        val z = event.values[2]

        // Calcular ángulos de inclinación / Calculate tilt angles
        val pitch = atan2(x, sqrt(y * y + z * z)) * 180f / Math.PI.toFloat()
        val roll = atan2(y, sqrt(x * x + z * z)) * 180f / Math.PI.toFloat()

        // Evaluar estabilidad / Evaluate stability
        val maxTilt = maxOf(abs(pitch), abs(roll))
        
        val newStatus = when {
            maxTilt > CAPSIZED_ANGLE -> {
                if (capsizedTime == 0L) {
                    capsizedTime = System.currentTimeMillis()
                }
                StabilityStatus.CAPSIZED
            }
            maxTilt > DANGER_ANGLE -> StabilityStatus.DANGER
            maxTilt > WARNING_ANGLE -> StabilityStatus.WARNING
            else -> {
                capsizedTime = 0
                StabilityStatus.STABLE
            }
        }

        // Verificar tiempo de vuelco / Check capsized time
        if (capsizedTime != 0L && 
            System.currentTimeMillis() - capsizedTime > SUBMERGED_TIMEOUT) {
            handleSubmergence()
        }

        if (_stability.value != newStatus) {
            _stability.postValue(newStatus)
            playAlertSound(newStatus)
            
            // Notificar emergencia / Notify emergency
            when (newStatus) {
                StabilityStatus.CAPSIZED -> 
                    emergencyNotifier?.sendEmergencyAlert(EmergencyType.CAPSIZE)
                StabilityStatus.SUBMERGED -> 
                    emergencyNotifier?.sendEmergencyAlert(EmergencyType.FLOODING)
                else -> { /* No emergency notification needed */ }
            }
        }
    }

    private fun processGyroscopeData(event: SensorEvent) {
        // Velocidad de rotación para análisis adicional
        // Rotation speed for additional analysis
        val rotationX = event.values[0]
        val rotationY = event.values[1]
        val rotationZ = event.values[2]

        // Detectar movimientos bruscos / Detect sudden movements
        val totalRotation = sqrt(
            rotationX * rotationX +
            rotationY * rotationY +
            rotationZ * rotationZ
        )

        if (totalRotation > 3.0f) { // Umbral de rotación / Rotation threshold
            _stability.postValue(StabilityStatus.WARNING)
            playAlertSound(StabilityStatus.WARNING)
        }
    }

    private fun handleSubmergence() {
        _stability.postValue(StabilityStatus.SUBMERGED)
        playAlertSound(StabilityStatus.SUBMERGED)
    }

    private fun playAlertSound(status: StabilityStatus) {
        when (status) {
            StabilityStatus.WARNING -> {
                // Tono corto de advertencia / Short warning tone
                toneGenerator.startTone(ToneGenerator.TONE_CDMA_ALERT_CALL_GUARD, 200)
            }
            StabilityStatus.DANGER -> {
                // Tono largo de peligro / Long danger tone
                toneGenerator.startTone(ToneGenerator.TONE_CDMA_EMERGENCY_RINGBACK, 500)
            }
            StabilityStatus.CAPSIZED -> {
                // Tono de emergencia / Emergency tone
                toneGenerator.startTone(ToneGenerator.TONE_CDMA_ABBR_ALERT, 1000)
            }
            StabilityStatus.SUBMERGED -> {
                // Tono continuo de emergencia / Continuous emergency tone
                toneGenerator.startTone(ToneGenerator.TONE_CDMA_EMERGENCY_RINGBACK, 2000)
            }
            else -> { /* No sound for stable status */ }
        }
    }

    override fun onAccuracyChanged(sensor: Sensor?, accuracy: Int) {
        // No se requiere implementación / No implementation required
    }
}
