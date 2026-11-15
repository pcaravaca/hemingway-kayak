package com.hemingway.kayak.utils

import org.junit.Assert.*
import org.junit.Test

class MapUtilsTest {

    @Test
    fun `GeoPoint creation with valid coordinates`() {
        val geoPoint = org.osmdroid.util.GeoPoint(40.7128, -74.0060)
        
        assertEquals(40.7128, geoPoint.latitude, 0.0001)
        assertEquals(-74.0060, geoPoint.longitude, 0.0001)
    }

    @Test
    fun `GeoPoint latitude bounds`() {
        val validNorth = org.osmdroid.util.GeoPoint(90.0, 0.0)
        val validSouth = org.osmdroid.util.GeoPoint(-90.0, 0.0)
        val invalidNorth = org.osmdroid.util.GeoPoint(91.0, 0.0)
        
        assertEquals(90.0, validNorth.latitude, 0.0001)
        assertEquals(-90.0, validSouth.latitude, 0.0001)
    }

    @Test
    fun `GeoPoint longitude bounds`() {
        val validEast = org.osmdroid.util.GeoPoint(0.0, 180.0)
        val validWest = org.osmdroid.util.GeoPoint(0.0, -180.0)
        
        assertEquals(180.0, validEast.longitude, 0.0001)
        assertEquals(-180.0, validWest.longitude, 0.0001)
    }

    @Test
    fun `GeoPoint distance calculation`() {
        val nyc = org.osmdroid.util.GeoPoint(40.7128, -74.0060)
        val london = org.osmdroid.util.GeoPoint(51.5074, -0.1278)
        
        val distance = nyc.distanceToAsDouble(london)
        
        assertTrue(distance > 5500000)
        assertTrue(distance < 5600000)
    }

    @Test
    fun `GeoPoint to string format`() {
        val geoPoint = org.osmdroid.util.GeoPoint(40.7128, -74.0060)
        val str = geoPoint.toString()
        
        assertTrue(str.contains("40.71"))
        assertTrue(str.contains("-74.00"))
    }
}
