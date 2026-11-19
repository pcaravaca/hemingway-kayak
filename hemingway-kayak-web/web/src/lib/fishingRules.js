// Reglas de "cuándo pescar". Deben coincidir con WeatherRepository.calculateFavorabilityScore
// de la app Android: el mismo número sale en la web y en el teléfono.

export const MOON_LABELS = {
  NEW_MOON: 'Luna nueva',
  WAXING_CRESCENT: 'Creciente',
  FIRST_QUARTER: 'Cuarto creciente',
  WAXING_GIBBOUS: 'Gibosa creciente',
  FULL_MOON: 'Luna llena',
  WANING_GIBBOUS: 'Gibosa menguante',
  THIRD_QUARTER: 'Cuarto menguante',
  WANING_CRESCENT: 'Menguante'
}

/** Fase aproximada por día del mes (igual que la app Android). */
export function moonPhaseFor(date) {
  const m = date.getDate() % 28
  if (m === 0) return 'NEW_MOON'
  if (m <= 3) return 'WAXING_CRESCENT'
  if (m <= 6) return 'FIRST_QUARTER'
  if (m <= 13) return 'WAXING_GIBBOUS'
  if (m === 14) return 'FULL_MOON'
  if (m <= 20) return 'WANING_GIBBOUS'
  if (m <= 24) return 'THIRD_QUARTER'
  return 'WANING_CRESCENT'
}

/**
 * @param {object} w clima del backend: temperature, windSpeed, cloudCover, pressure
 * @param {string} phase una de las claves de MOON_LABELS
 * @param {string} pressureTrend STABLE | RISING | FALLING (sin histórico: STABLE)
 */
export function favorabilityScore(w, phase, pressureTrend = 'STABLE') {
  let score = 50

  const temp = w.temperature
  if (temp >= 18 && temp <= 24) score += 15
  else if (temp >= 15 && temp <= 28) score += 10
  else if (temp >= 10 && temp <= 30) score += 5
  else if (temp < 5 || temp > 35) score -= 10

  const wind = w.windSpeed ?? 0
  if (wind < 2) score += 10
  else if (wind <= 5) score += 5
  else if (wind <= 10) score -= 5
  else score -= 15

  if (pressureTrend === 'STABLE') score += 10
  else if (pressureTrend === 'RISING') score += 5
  else if (pressureTrend === 'FALLING') score -= 5

  if (phase === 'FULL_MOON' || phase === 'NEW_MOON') score += 15
  else if (phase === 'FIRST_QUARTER' || phase === 'THIRD_QUARTER') score += 10

  return Math.max(0, Math.min(100, score))
}

export function scoreLabel(score) {
  if (score >= 70) return 'Excelente para salir'
  if (score >= 50) return 'Buenas condiciones'
  if (score >= 35) return 'Regular, se puede'
  return 'Mejor esperar'
}

export function suggestedTechniques(w, score) {
  const list = []
  if (w.windSpeed > 8) list.push('Casting a sotavento')
  if (w.windSpeed < 3) list.push('Técnicas de superficie')
  if ((w.cloudCover ?? 0) > 80) list.push('Señuelos oscuros')
  if ((w.cloudCover ?? 0) < 30) list.push('Señuelos brillantes')
  if (score > 70) list.push('Aguas superficiales', 'Cerca de vegetación')
  if (score < 40) list.push('Aguas profundas', 'Pesca lenta')
  return list
}

export function suggestedBaits(w, score) {
  const list = []
  if (w.temperature < 15) list.push('Jigs')
  if (w.temperature > 25) list.push('Topwater')
  if (score > 70) list.push('Crankbaits', 'Spinnerbaits')
  if (score < 40) list.push('Soft plastics', 'Carolina rig')
  return list
}

/** Mejor momento para salir, con la misma lógica que la app. */
export function bestTimeToFish(now, phase) {
  const hour = now.getHours()
  const at = (hourOfDay, addDays = 0) => {
    const t = new Date(now)
    if (addDays) t.setDate(t.getDate() + addDays)
    t.setHours(hourOfDay, 0, 0, 0)
    return t
  }

  if (hour >= 18) return at(6, 1)
  if (phase === 'FULL_MOON' || phase === 'NEW_MOON') {
    return hour < 12 ? at(18) : at(6, 1)
  }
  if (hour < 6) return now
  return at(18)
}
