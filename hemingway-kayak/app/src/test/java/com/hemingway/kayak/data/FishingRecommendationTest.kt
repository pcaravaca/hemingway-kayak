package com.hemingway.kayak.data

import org.junit.Assert.*
import org.junit.Test
import java.util.Date

class FishingRecommendationTest {

    @Test
    fun `FishingRecommendation creation with all fields`() {
        val bestTime = Date()
        
        val recommendation = FishingRecommendation(
            favorabilityScore = 85,
            pressureTrend = PressureTrend.RISING,
            moonPhase = MoonPhase.FULL_MOON,
            bestTimeToFish = bestTime,
            suggestedTechniques = listOf("Casting", "Topwater"),
            suggestedBaits = listOf("Crankbait", "Spinnerbait")
        )

        assertEquals(85, recommendation.favorabilityScore)
        assertEquals(PressureTrend.RISING, recommendation.pressureTrend)
        assertEquals(MoonPhase.FULL_MOON, recommendation.moonPhase)
        assertEquals(bestTime, recommendation.bestTimeToFish)
        assertEquals(2, recommendation.suggestedTechniques.size)
        assertEquals(2, recommendation.suggestedBaits.size)
    }

    @Test
    fun `FishingRecommendation with null best time`() {
        val recommendation = FishingRecommendation(
            favorabilityScore = 30,
            pressureTrend = PressureTrend.FALLING,
            moonPhase = MoonPhase.NEW_MOON,
            bestTimeToFish = null,
            suggestedTechniques = emptyList(),
            suggestedBaits = emptyList()
        )

        assertNull(recommendation.bestTimeToFish)
        assertTrue(recommendation.suggestedTechniques.isEmpty())
    }

    @Test
    fun `PressureTrend enum values`() {
        assertEquals(3, PressureTrend.values().size)
        assertNotNull(PressureTrend.RISING)
        assertNotNull(PressureTrend.FALLING)
        assertNotNull(PressureTrend.STABLE)
    }

    @Test
    fun `MoonPhase enum values`() {
        assertEquals(8, MoonPhase.values().size)
        assertNotNull(MoonPhase.NEW_MOON)
        assertNotNull(MoonPhase.FULL_MOON)
        assertNotNull(MoonPhase.WAXING_CRESCENT)
        assertNotNull(MoonPhase.WANING_CRESCENT)
    }
}
