package com.hemingway.kayak.data.remote

import android.content.Context
import com.hemingway.kayak.BuildConfig
import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import dagger.hilt.android.qualifiers.ApplicationContext
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Arma el cliente de la API con el JWT ya guardado. Cada request protegido sale
 * con `Authorization: Bearer <token>`; si no hay token, la API responde 401.
 */
@Singleton
class BackendClient @Inject constructor(
    @ApplicationContext context: Context,
    private val tokenStore: TokenStore,
    private val okHttpClient: OkHttpClient,
    private val moshi: Moshi
) {
    private val authInterceptor = Interceptor { chain ->
        val request = chain.request().newBuilder()
        tokenStore.token?.let { request.header("Authorization", "Bearer $it") }
        chain.proceed(request.build())
    }

    val api: KayakApiService by lazy {
        Retrofit.Builder()
            .baseUrl(BuildConfig.API_BASE_URL)
            .client(okHttpClient.newBuilder().addInterceptor(authInterceptor).build())
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .build()
            .create(KayakApiService::class.java)
    }

    // La app se conecta a la URL del server (Tailscale o LAN), no a la del emulador.
    companion object {
        fun normalize(baseUrl: String): String = if (baseUrl.endsWith("/")) baseUrl else "$baseUrl/"
    }
}
