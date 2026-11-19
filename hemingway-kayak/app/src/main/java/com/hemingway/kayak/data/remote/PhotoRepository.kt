package com.hemingway.kayak.data.remote

import android.content.Context
import android.net.Uri
import dagger.hilt.android.qualifiers.ApplicationContext
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.toRequestBody
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Sube las fotos de captura. El backend valida el contenido (no el header), 8 MB,
 * JPEG/PNG/WebP/HEIC. Devuelve la URL pública (/uploads/...) o null.
 */
@Singleton
class PhotoRepository @Inject constructor(
    @ApplicationContext private val context: Context,
    private val client: BackendClient
) {
    suspend fun upload(uri: Uri): String? = runCatching {
        val resolver = context.contentResolver
        val bytes = resolver.openInputStream(uri)?.use { it.readBytes() }
            ?: throw IllegalStateException("No se pudo leer la foto")
        val name = uri.lastPathSegment?.substringAfterLast('/') ?: "captura.jpg"
        val part = MultipartBody.Part.createFormData(
            "photo",
            name,
            bytes.toRequestBody("image/*".toMediaTypeOrNull())
        )
        client.api.uploadPhoto(part).url
    }.getOrNull()
}
