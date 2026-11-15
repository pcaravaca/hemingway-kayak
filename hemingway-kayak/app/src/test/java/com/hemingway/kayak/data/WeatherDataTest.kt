package com.hemingway.kayak.data

import org.junit.Assert.*
import org.junit.Test
import java.util.Date

class WeatherDataTest {

    @Test
    fun `WeatherData creation with required fields`() {
        val now = Date()
        val expiration = Date(now.time + 3600000)

        val weatherData = WeatherData(
            latitude = 40.7128,
            longitude = -74.0060,
            locationName = "New York",
            temperature = 22.5f,
            feelsLike = 21.0f,
            pressure = 1013.0f,
            humidity = 65,
            windSpeed = 5.5f,
            windDirection = 270,
            cloudCover = 40,
            precipitation = 0.0f,
            timestamp = now,
            lastUpdated = now,
            expirationTime = expiration
        )

        assertEquals(40.7128, weatherData.latitude, 0.0001)
        assertEquals(-74.0060, weatherData.longitude, 0.0001)
        assertEquals("New York", weatherData.locationName)
        assertEquals(22.5f, weatherData.temperature, 0.1f)
        assertEquals(21.0f, weatherData.feelsLike, 0.1f)
        assertEquals(1013.0f, weatherData.pressure, 0.1f)
        assertEquals(65, weatherData.humidity)
        assertEquals(5.5f, weatherData.windSpeed, 0.1f)
        assertEquals(270, weatherData.windDirection)
        assertEquals(40, weatherData.cloudCover)
        assertEquals(0.0f, weatherData.precipitation, 0.1f)
    }

    @Test
    fun `WeatherData default description is empty`() {
        val now = Date()
        val expiration = Date(now.time + 3600000)

        val weatherData = WeatherData(
            latitude = 40.7128,
            longitude = -74.0060,
            locationName = "NYC",
            temperature = 20f,
            feelsLike = 19f,
            pressure = 1015f,
            humidity = 70,
            windSpeed = 3f,
            windDirection = 90,
            cloudCover = 80,
            precipitation = 0f,
            timestamp = now,
            lastUpdated = now,
            expirationTime = expiration
        )

        assertEquals("", weatherData.description)
    }
}
