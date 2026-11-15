package com.hemingway.kayak.data

import org.junit.Assert.*
import org.junit.Test

class ResourceTest {

    @Test
    fun `ResourceSuccess contains data`() {
        val weatherData = WeatherData(
            latitude = 40.7128,
            longitude = -74.0060,
            locationName = "NYC",
            temperature = 22f,
            feelsLike = 21f,
            pressure = 1013f,
            humidity = 65,
            windSpeed = 5f,
            windDirection = 180,
            cloudCover = 50,
            precipitation = 0f,
            timestamp = java.util.Date(),
            lastUpdated = java.util.Date(),
            expirationTime = java.util.Date()
        )
        
        val resource = Resource.Success(weatherData, isFromCache = false)
        
        assertTrue(resource.data != null)
        assertEquals(weatherData, resource.data)
        assertNull(resource.message)
        assertFalse(resource.isFromCache)
    }

    @Test
    fun `ResourceSuccess with cached data`() {
        val resource = Resource.Success<String>("cached value", isFromCache = true)
        
        assertEquals("cached value", resource.data)
        assertTrue(resource.isFromCache)
    }

    @Test
    fun `ResourceError contains message and data`() {
        val weatherData = WeatherData(
            latitude = 0.0,
            longitude = 0.0,
            locationName = "",
            temperature = 0f,
            feelsLike = 0f,
            pressure = 0f,
            humidity = 0,
            windSpeed = 0f,
            windDirection = 0,
            cloudCover = 0,
            precipitation = 0f,
            timestamp = java.util.Date(),
            lastUpdated = java.util.Date(),
            expirationTime = java.util.Date()
        )
        
        val resource = Resource.Error("Network error", weatherData)
        
        assertEquals("Network error", resource.message)
        assertNotNull(resource.data)
        assertEquals(weatherData, resource.data)
    }

    @Test
    fun `ResourceLoading contains data`() {
        val resource = Resource.Loading<String>("some data")
        
        assertEquals("some data", resource.data)
    }
}
