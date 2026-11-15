package com.hemingway.kayak.data

import androidx.room.*
import kotlinx.coroutines.flow.Flow

/**
 * DAO para operaciones de base de datos con modelos de kayak.
 * DAO for database operations with kayak models.
 */
@Dao
interface KayakDao {
    /**
     * Insertar un nuevo modelo de kayak.
     * Insert a new kayak model.
     */
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertKayakModel(kayakModel: KayakModel)

    /**
     * Obtener todos los modelos de kayak guardados.
     * Get all saved kayak models.
     */
    @Query("SELECT * FROM kayak_models")
    fun getAllKayakModels(): Flow<List<KayakModel>>

    /**
     * Obtener un modelo de kayak por su ID.
     * Get a kayak model by its ID.
     */
    @Query("SELECT * FROM kayak_models WHERE id = :kayakId")
    fun getKayakModelById(kayakId: Int): Flow<KayakModel>

    /**
     * Actualizar un modelo de kayak existente.
     * Update an existing kayak model.
     */
    @Update
    suspend fun updateKayakModel(kayakModel: KayakModel)

    /**
     * Eliminar un modelo de kayak.
     * Delete a kayak model.
     */
    @Delete
    suspend fun deleteKayakModel(kayakModel: KayakModel)

    /**
     * Buscar modelos de kayak por tipo.
     * Search kayak models by type.
     */
    @Query("SELECT * FROM kayak_models WHERE type = :kayakType")
    fun getKayakModelsByType(kayakType: KayakType): Flow<List<KayakModel>>
}
