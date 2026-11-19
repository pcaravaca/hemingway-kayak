import { describe, it, expect } from 'vitest'
import {
  MOON_LABELS,
  bestTimeToFish,
  favorabilityScore,
  moonPhaseFor,
  scoreLabel,
  suggestedBaits,
  suggestedTechniques
} from '@/lib/fishingRules'

// Las reglas deben coincidir con WeatherRepository de Android. Si cambia una,
// cambian las dos: la web y el teléfono tienen que dar la misma nota.

describe('moonPhaseFor', () => {
  const phase = (day) => moonPhaseFor(new Date(2026, 10, day))

  it('día 14 = luna llena', () => {
    expect(phase(14)).toBe('FULL_MOON')
  })

  it('día 28 = luna nueva (28 % 28 = 0)', () => {
    expect(phase(28)).toBe('NEW_MOON')
  })

  it('día 5 = cuarto creciente', () => {
    expect(phase(5)).toBe('FIRST_QUARTER')
  })

  it('todos los valores tienen etiqueta en español', () => {
    for (let day = 1; day <= 30; day++) {
      expect(MOON_LABELS[phase(day)], `día ${day}`).toBeTruthy()
    }
  })
})

describe('favorabilityScore', () => {
  const base = { temperature: 21, windSpeed: 1.5, cloudCover: 20, pressure: 1013 }

  it('condiciones ideales (21 °C, viento 1.5, luna llena) = 100', () => {
    // 50 + 15 (temp) + 10 (viento) + 10 (presión estable) + 15 (luna llena)
    expect(favorabilityScore(base, 'FULL_MOON')).toBe(100)
  })

  it('calor extremo y viento fuerte baja la nota', () => {
    const malo = { temperature: 38, windSpeed: 15, cloudCover: 50, pressure: 1013 }
    // 50 - 10 (temp) - 15 (viento) + 10 (estable) + 0 (sin luna)
    expect(favorabilityScore(malo, 'WANING_CRESCENT')).toBe(35)
  })

  it('nunca se pasa de 0 ni de 100', () => {
    const extremo = { temperature: -20, windSpeed: 40, cloudCover: 100, pressure: 1013 }
    const score = favorabilityScore(extremo, 'WANING_CRESCENT', 'FALLING')
    expect(score).toBeGreaterThanOrEqual(0)
    expect(score).toBeLessThanOrEqual(100)
  })

  it('la tendencia de presión suma o resta', () => {
    const estable = favorabilityScore(base, 'WANING_CRESCENT', 'STABLE')
    const subiendo = favorabilityScore(base, 'WANING_CRESCENT', 'RISING')
    const bajando = favorabilityScore(base, 'WANING_CRESCENT', 'FALLING')
    expect(subiendo).toBe(estable - 5)
    expect(bajando).toBe(estable - 15)
  })

  it('sin viento reportado no rompe', () => {
    expect(favorabilityScore({ temperature: 21 }, 'FULL_MOON')).toBe(100)
  })
})

describe('scoreLabel', () => {
  it('los cortes', () => {
    expect(scoreLabel(85)).toBe('Excelente para salir')
    expect(scoreLabel(60)).toBe('Buenas condiciones')
    expect(scoreLabel(40)).toBe('Regular, se puede')
    expect(scoreLabel(10)).toBe('Mejor esperar')
  })
})

describe('suggestedTechniques', () => {
  it('viento fuerte pide casting a sotavento', () => {
    expect(suggestedTechniques({ windSpeed: 10, cloudCover: 50 }, 60)).toContain('Casting a sotavento')
  })

  it('mucha nube pide señuelos oscuros', () => {
    expect(suggestedTechniques({ windSpeed: 4, cloudCover: 90 }, 60)).toContain('Señuelos oscuros')
  })

  it('nota alta agrega aguas superficiales', () => {
    expect(suggestedTechniques({ windSpeed: 4, cloudCover: 50 }, 80)).toContain('Aguas superficiales')
  })

  it('nota baja agrega pesca lenta', () => {
    expect(suggestedTechniques({ windSpeed: 4, cloudCover: 50 }, 20)).toContain('Pesca lenta')
  })
})

describe('suggestedBaits', () => {
  it('frío → jigs', () => {
    expect(suggestedBaits({ temperature: 12 }, 60)).toContain('Jigs')
  })

  it('calor → topwater', () => {
    expect(suggestedBaits({ temperature: 28 }, 60)).toContain('Topwater')
  })
})

describe('bestTimeToFish', () => {
  it('después de las 18h sugiere mañana temprano', () => {
    const result = bestTimeToFish(new Date(2026, 10, 19, 20, 0), 'WANING_CRESCENT')
    expect(result.getHours()).toBe(6)
    expect(result.getDate()).toBe(20)
  })

  it('luna llena antes del mediodía sugiere el atardecer', () => {
    const result = bestTimeToFish(new Date(2026, 10, 19, 9, 0), 'FULL_MOON')
    expect(result.getHours()).toBe(18)
    expect(result.getDate()).toBe(19)
  })
})
