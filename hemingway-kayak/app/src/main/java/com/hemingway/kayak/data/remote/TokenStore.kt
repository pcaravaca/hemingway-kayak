package com.hemingway.kayak.data.remote

import android.content.Context
import androidx.core.content.edit

/**
 * Guarda el JWT que devuelve /api/v1/auth/login. Sin él ninguna llamada
 * protegida (tracking, torneos, usuarios) funciona.
 */
class TokenStore(context: Context) {

    private val prefs = context.getSharedPreferences("hk_auth", Context.MODE_PRIVATE)

    var token: String?
        get() = prefs.getString(KEY_TOKEN, null)
        set(value) = prefs.edit { if (value.isNullOrBlank()) remove(KEY_TOKEN) else putString(KEY_TOKEN, value) }

    var userId: String?
        get() = prefs.getString(KEY_USER_ID, null)
        set(value) = prefs.edit { if (value.isNullOrBlank()) remove(KEY_USER_ID) else putString(KEY_USER_ID, value) }

    val isLoggedIn: Boolean get() = !token.isNullOrBlank()

    fun clear() {
        prefs.edit { clear() }
    }

    private companion object {
        const val KEY_TOKEN = "token"
        const val KEY_USER_ID = "userId"
    }
}
