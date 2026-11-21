package com.hemingway.kayak.data.remote

import javax.inject.Inject
import javax.inject.Singleton

/** Kayaks del usuario: listado, alta, edicion y borrado. */
@Singleton
class KayakRepository @Inject constructor(
    private val client: BackendClient
) {
    suspend fun all(): List<KayakDto> = try {
        client.api.kayaks()
    } catch (e: Exception) {
        emptyList()
    }

    suspend fun create(request: KayakRequest): Boolean = try {
        client.api.createKayak(request)
        true
    } catch (e: Exception) {
        false
    }

    suspend fun update(id: String, request: KayakRequest): Boolean = try {
        client.api.updateKayak(id, request)
        true
    } catch (e: Exception) {
        false
    }

    suspend fun delete(id: String): Boolean = try {
        client.api.deleteKayak(id)
        true
    } catch (e: Exception) {
        false
    }
}
