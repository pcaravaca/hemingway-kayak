package com.hemingway.kayak.data

import androidx.room.*
import kotlinx.coroutines.flow.Flow

/**
 * DAO para operaciones de base de datos con puntos de pesca.
 * DAO for database operations with fishing spots.
 */
@Dao
interface FishingSpotDao {
    /**
     * Insertar un nuevo punto de pesca.
     * Insert a new fishing spot.
     */
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertFishingSpot(fishingSpot: FishingSpot): Long
    
    /**
     * Actualizar un punto de pesca existente.
     * Update an existing fishing spot.
     */
    @Update
    suspend fun updateFishingSpot(fishingSpot: FishingSpot)
    
    /**
     * Eliminar un punto de pesca.
     * Delete a fishing spot.
     */
    @Delete
    suspend fun deleteFishingSpot(fishingSpot: FishingSpot)
    
    /**
     * Obtener todos los puntos de pesca.
     * Get all fishing spots.
     */
    @Query("SELECT * FROM fishing_spots ORDER BY createdAt DESC")
    fun getAllFishingSpots(): Flow<List<FishingSpot>>
    
    /**
     * Obtener un punto de pesca por ID.
     * Get a fishing spot by ID.
     */
    @Query("SELECT * FROM fishing_spots WHERE id = :id")
    fun getFishingSpotById(id: Int): Flow<FishingSpot?>
    
    /**
     * Buscar puntos de pesca cercanos a una ubicación.
     * Search for fishing spots near a location.
     */
    @Query("""
        SELECT * FROM fishing_spots 
        WHERE (latitude BETWEEN :minLat AND :maxLat) 
        AND (longitude BETWEEN :minLng AND :maxLng)
        ORDER BY createdAt DESC
    """)
    fun getNearbyFishingSpots(minLat: Double, maxLat: Double, minLng: Double, maxLng: Double): Flow<List<FishingSpot>>
    
    /**
     * Obtener los puntos de pesca con mayor puntuación.
     * Get the highest rated fishing spots.
     */
    @Query("SELECT * FROM fishing_spots ORDER BY rating DESC LIMIT :limit")
    fun getTopRatedFishingSpots(limit: Int = 10): Flow<List<FishingSpot>>
    
    /**
     * Obtener los puntos de pesca compartidos.
     * Get shared fishing spots.
     */
    @Query("SELECT * FROM fishing_spots WHERE isShared = 1 ORDER BY createdAt DESC")
    fun getSharedFishingSpots(): Flow<List<FishingSpot>>
}
