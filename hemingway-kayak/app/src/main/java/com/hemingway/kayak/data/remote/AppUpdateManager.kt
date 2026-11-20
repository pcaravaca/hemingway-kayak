package com.hemingway.kayak.data.remote

import android.content.Context
import android.content.pm.PackageManager
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Busca una versión nueva de la app al arrancar. Compara el versionCode del APK
 * instalado contra el `version.json` que publica el servidor junto con el APK.
 * Si el archivo no está o no se puede leer, no pasa nada: no hay aviso.
 */
@Singleton
class AppUpdateManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)

    private val _update = MutableStateFlow<UpdateInfo?>(null)
    val update: StateFlow<UpdateInfo?> = _update.asStateFlow()

    data class UpdateInfo(
        val versionCode: Int,
        val versionName: String,
        val apkUrl: String,
        val releaseNotes: String?
    )

    /** Una consulta por proceso: si el servidor no responde, no se reintenta. */
    private var checked = false

    fun checkForUpdate() {
        if (checked) return
        checked = true
        scope.launch {
            val info = fetchLatest() ?: return@launch
            if (info.versionCode > currentVersionCode()) {
                _update.value = info
            }
        }
    }

    fun dismiss() {
        _update.value = null
    }

    @Suppress("DEPRECATION")
    private fun currentVersionCode(): Long = try {
        context.packageManager.getPackageInfo(context.packageName, 0).versionCode.toLong()
    } catch (e: Exception) {
        0L
    }

    private fun fetchLatest(): UpdateInfo? {
        return try {
            val url = URL(KayakApiService.DEFAULT_DOWNLOAD_BASE_URL + "version.json")
            val connection = (url.openConnection() as HttpURLConnection).apply {
                connectTimeout = 5000
                readTimeout = 5000
            }
            try {
                if (connection.responseCode != HttpURLConnection.HTTP_OK) return null
                val body = connection.inputStream.bufferedReader().use { it.readText() }
                val json = JSONObject(body)
                val code = json.getInt("versionCode")
                UpdateInfo(
                    versionCode = code,
                    versionName = json.optString("versionName", code.toString()),
                    apkUrl = json.optString(
                        "apkUrl",
                        KayakApiService.DEFAULT_DOWNLOAD_BASE_URL + "HemingwayKayak.apk"
                    ),
                    releaseNotes = json.optString("releaseNotes", "").ifEmpty { null }
                )
            } finally {
                connection.disconnect()
            }
        } catch (e: Exception) {
            null
        }
    }
}
