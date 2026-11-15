package com.hemingway.kayak.data

import org.junit.Assert.*
import org.junit.Test
import java.util.Date

class FishingSpotTest {

    @Test
    fun `FishingSpot creation with default values`() {
        val spot = FishingSpot(
            name = "Test Spot",
            description = "Test Description",
            latitude = 40.7128,
            longitude = -74.0060,
            createdAt = Date()
        )

        assertEquals("Test Spot", spot.name)
        assertEquals("Test Description", spot.description)
        assertEquals(40.7128, spot.latitude, 0.0001)
        assertEquals(-74.0060, spot.longitude, 0.0001)
        assertEquals(0, spot.rating)
        assertEquals(0, spot.fishCaught)
        assertFalse(spot.isShared)
    }

    @Test
    fun `FishingSpot creation with all values`() {
        val createdAt = Date()
        val lastVisited = Date()
        
        val spot = FishingSpot(
            id = 1,
            name = "Hobie Spot",
            description = "Great for kayaks",
            latitude = 25.7617,
            longitude = -80.1918,
            createdAt = createdAt,
            lastVisited = lastVisited,
            rating = 5,
            notes = "Best at dawn",
            fishCaught = 10,
            isShared = true,
            imageUri = "file:///image.jpg"
        )

        assertEquals(1, spot.id)
        assertEquals(5, spot.rating)
        assertEquals(10, spot.fishCaught)
        assertTrue(spot.isShared)
        assertEquals("file:///image.jpg", spot.imageUri)
        assertEquals(lastVisited, spot.lastVisited)
    }
}
