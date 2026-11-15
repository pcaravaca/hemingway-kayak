package com.hemingway.kayak.data

import androidx.room.Entity
import androidx.room.PrimaryKey
import java.util.Date

/**
 * Modelo de datos para un punto de pesca.
 * Data model for a fishing spot.
 */
@Entity(tableName = "fishing_spots")
data class FishingSpot(
    @PrimaryKey(autoGenerate = true) val id: Int = 0,
    val name: String,
    val description: String,
    val latitude: Double,
    val longitude: Double,
    val createdAt: Date,
    val lastVisited: Date? = null,
    val rating: Int = 0, // 0-5 estrellas / 0-5 stars
    val notes: String = "",
    val fishCaught: Int = 0,
    val isShared: Boolean = false,
    val imageUri: String? = null
)
