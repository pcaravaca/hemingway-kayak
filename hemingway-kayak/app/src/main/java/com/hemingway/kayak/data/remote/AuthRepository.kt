package com.hemingway.kayak.data.remote

import javax.inject.Inject
import javax.inject.Singleton

/** Login/registro contra el backend. Guarda el token si sale bien. */
@Singleton
class AuthRepository @Inject constructor(
    private val client: BackendClient,
    private val tokenStore: TokenStore
) {
    val isLoggedIn: Boolean get() = tokenStore.isLoggedIn

    /** Devuelve null si salió bien; si no, el mensaje de error del backend. */
    suspend fun login(email: String, password: String): String? = try {
        val response = client.api.login(AuthRequest(email.trim().lowercase(), password.trim()))
        tokenStore.token = response.token
        tokenStore.userId = response.user.id
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo entrar"
    }

    suspend fun register(name: String, email: String, password: String): String? = try {
        val response = client.api.register(
            RegisterRequest(
                name = name,
                email = email.trim().lowercase(),
                username = email.substringBefore("@"),
                password = password
            )
        )
        tokenStore.token = response.token
        tokenStore.userId = response.user.id
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo crear la cuenta"
    }

    fun logout() = tokenStore.clear()
}
