package com.hemingway.kayak.di

import android.content.Context
import com.hemingway.kayak.data.FishingSpotDao
import com.hemingway.kayak.data.KayakDao
import com.hemingway.kayak.data.KayakDatabase
import com.hemingway.kayak.data.WeatherApiService
import com.hemingway.kayak.data.WeatherDao
import com.hemingway.kayak.data.WeatherRepository
import com.hemingway.kayak.data.remote.TokenStore
import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import okhttp3.Cache
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DataModule {

    @Singleton
    @Provides
    fun provideDatabase(@ApplicationContext context: Context): KayakDatabase {
        return KayakDatabase.getDatabase(context)
    }

    @Singleton
    @Provides
    fun provideKayakDao(database: KayakDatabase): KayakDao {
        return database.kayakDao()
    }

    @Singleton
    @Provides
    fun provideWeatherDao(database: KayakDatabase): WeatherDao {
        return database.weatherDao()
    }

    @Singleton
    @Provides
    fun provideFishingSpotDao(database: KayakDatabase): FishingSpotDao {
        return database.fishingSpotDao()
    }

    @Singleton
    @Provides
    fun provideOkHttpClient(@ApplicationContext context: Context): OkHttpClient {
        val cacheSize = 10 * 1024 * 1024L
        val cache = Cache(context.cacheDir, cacheSize)

        return OkHttpClient.Builder()
            .cache(cache)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .writeTimeout(30, TimeUnit.SECONDS)
            .build()
    }

    @Singleton
    @Provides
    fun provideMoshi(): Moshi {
        return Moshi.Builder()
            .add(KotlinJsonAdapterFactory())
            .build()
    }

    // JWT de la sesión contra el backend de Hemingway Kayak.
    @Singleton
    @Provides
    fun provideTokenStore(@ApplicationContext context: Context): TokenStore {
        return TokenStore(context)
    }

    @Singleton
    @Provides
    fun provideWeatherApiService(okHttpClient: OkHttpClient, moshi: Moshi): WeatherApiService {
        return Retrofit.Builder()
            .baseUrl(WeatherApiService.BASE_URL)
            .client(okHttpClient)
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .build()
            .create(WeatherApiService::class.java)
    }

    @Singleton
    @Provides
    fun provideWeatherRepository(
        weatherDao: WeatherDao,
        weatherApiService: WeatherApiService
    ): WeatherRepository {
        return WeatherRepository(weatherDao, weatherApiService)
    }
}
