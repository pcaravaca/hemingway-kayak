package com.hemingway.kayak.emergency

import org.junit.Assert.*
import org.junit.Test

class EmergencyNotifierTest {

    @Test
    fun `EmergencyState enum values`() {
        val states = EmergencyNotifier.EmergencyState.values()
        
        assertEquals(5, states.size)
        assertTrue(states.contains(EmergencyNotifier.EmergencyState.IDLE))
        assertTrue(states.contains(EmergencyNotifier.EmergencyState.TRIGGERED))
        assertTrue(states.contains(EmergencyNotifier.EmergencyState.SENDING))
        assertTrue(states.contains(EmergencyNotifier.EmergencyState.SENT))
        assertTrue(states.contains(EmergencyNotifier.EmergencyState.FAILED))
    }

    @Test
    fun `EmergencyInfo creation`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "CAPSIZE"
        )
        
        assertEquals(40.7128, info.latitude, 0.0001)
        assertEquals(-74.0060, info.longitude, 0.0001)
        assertEquals("CAPSIZE", info.type)
        assertEquals(0f, info.driftDistance, 0.001f)
        assertEquals(0f, info.driftDirection, 0.001f)
    }

    @Test
    fun `EmergencyInfo compass direction North`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 0f
        )
        
        assertEquals("N", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction East`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 90f
        )
        
        assertEquals("E", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction South`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 180f
        )
        
        assertEquals("S", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction West`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 270f
        )
        
        assertEquals("W", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction NE`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 45f
        )
        
        assertEquals("NE", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction SW`() {
        val info = EmergencyNotifier.EmergencyInfo(
            latitude = 40.7128,
            longitude = -74.0060,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 225f
        )
        
        assertEquals("SW", info.compassDirection)
    }

    @Test
    fun `EmergencyInfo compass direction at boundary N`() {
        val infoNorth = EmergencyNotifier.EmergencyInfo(
            latitude = 0.0, longitude = 0.0,
            timestamp = System.currentTimeMillis(),
            type = "TEST",
            driftDirection = 337.5f
        )
        assertEquals("N", infoNorth.compassDirection)
    }
}
