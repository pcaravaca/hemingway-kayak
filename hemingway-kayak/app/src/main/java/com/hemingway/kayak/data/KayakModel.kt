package com.hemingway.kayak.data

import androidx.room.Entity
import androidx.room.PrimaryKey

/**
 * Representa un modelo de kayak de pesca con sus especificaciones técnicas.
 * Fishing kayak model representation with technical specifications.
 */
@Entity(tableName = "kayak_models")
data class KayakModel(
    @PrimaryKey(autoGenerate = true) val id: Int = 0,
    val name: String, // Nombre del kayak / Kayak name
    val type: KayakType, // Tipo de kayak / Kayak type
    val length: Float, // Longitud en metros / Length in meters
    val width: Float, // Ancho en metros / Width in meters
    val weight: Float, // Peso del kayak en kg / Kayak weight in kg
    val maxCapacity: Float, // Capacidad máxima de carga en kg / Maximum load capacity in kg
    val material: KayakMaterial, // Material de construcción / Construction material
    val propulsionType: PropulsionType // Tipo de propulsión / Propulsion type
)

/**
 * Tipos de kayaks de pesca más comunes.
 * Most common fishing kayak types.
 */
enum class KayakType {
    SIT_ON_TOP,  // Sentado sobre el kayak / Seated on top
    SIT_INSIDE,  // Sentado dentro del kayak / Seated inside
    INFLATABLE,  // Kayak inflable / Inflatable kayak
    PEDAL_DRIVE  // Kayak con impulso de pedales / Pedal-driven kayak
}

/**
 * Materiales de construcción de kayaks.
 * Kayak construction materials.
 */
enum class KayakMaterial {
    POLYETHYLENE,  // Polietileno rotomoldeado / Rotomolded polyethylene
    FIBERGLASS,    // Fibra de vidrio / Fiberglass
    CARBON_FIBER,  // Fibra de carbono / Carbon fiber
    INFLATABLE_PVC // PVC para kayaks inflables / PVC for inflatable kayaks
}

/**
 * Tipos de propulsión para kayaks.
 * Kayak propulsion types.
 */
enum class PropulsionType {
    PADDLE,       // Remo / Paddle
    PEDAL_DRIVE,  // Impulso de pedales / Pedal drive
    MOTOR        // Motor eléctrico / Electric motor
}

/**
 * Datos predeterminados de modelos de kayaks populares.
 * Default data for popular kayak models.
 */
object DefaultKayakModels {
    val defaultModels = listOf(
        KayakModel(
            name = "Hobie Mirage Pro Angler 12",
            type = KayakType.PEDAL_DRIVE,
            length = 3.66f,  // 12 pies
            width = 0.91f,   // 36 pulgadas
            weight = 45.5f,  // 100 libras
            maxCapacity = 204.1f, // 450 libras
            material = KayakMaterial.POLYETHYLENE,
            propulsionType = PropulsionType.PEDAL_DRIVE
        ),
        KayakModel(
            name = "Wilderness Systems Radar 135",
            type = KayakType.SIT_ON_TOP,
            length = 4.11f,  // 13.5 pies
            width = 0.86f,   // 34 pulgadas
            weight = 40.8f,  // 90 libras
            maxCapacity = 181.4f, // 400 libras
            material = KayakMaterial.POLYETHYLENE,
            propulsionType = PropulsionType.PADDLE
        ),
        KayakModel(
            name = "Old Town Sportsman BigWater PDL 132",
            type = KayakType.PEDAL_DRIVE,
            length = 4.02f,  // 13.2 pies
            width = 0.91f,   // 36 pulgadas
            weight = 49.9f,  // 110 libras
            maxCapacity = 226.8f, // 500 libras
            material = KayakMaterial.POLYETHYLENE,
            propulsionType = PropulsionType.PEDAL_DRIVE
        )
    )
}

/**
 * Utilidad para calcular rendimiento y características del kayak.
 * Utility for calculating kayak performance and characteristics.
 */
object KayakPerformanceCalculator {
    /**
     * Calcula la velocidad máxima estimada del kayak.
     * Calculates estimated maximum kayak speed.
     */
    fun calculateMaxSpeed(kayak: KayakModel, paddlerWeight: Float): Float {
        // Fórmula simplificada de velocidad / Simplified speed formula
        return when (kayak.propulsionType) {
            PropulsionType.PADDLE -> (0.5f * kayak.length / (kayak.weight + paddlerWeight)) * 3.6f
            PropulsionType.PEDAL_DRIVE -> (0.7f * kayak.length / (kayak.weight + paddlerWeight)) * 3.6f
            PropulsionType.MOTOR -> 10f // Velocidad base para kayaks con motor / Base speed for motor kayaks
        }
    }

    /**
     * Calcula la estabilidad del kayak.
     * Calculates kayak stability.
     */
    fun calculateStability(kayak: KayakModel): Float {
        return kayak.width * 100 / kayak.length
    }

    /**
     * Calcula la eficiencia de remada o pedaleo.
     * Calculates paddling or pedaling efficiency.
     */
    fun calculateEfficiency(kayak: KayakModel): Float {
        return when (kayak.propulsionType) {
            PropulsionType.PADDLE -> 0.6f
            PropulsionType.PEDAL_DRIVE -> 0.8f
            PropulsionType.MOTOR -> 0.9f
        }
    }
}
