<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import {
  MOON_LABELS,
  bestTimeToFish,
  favorabilityScore,
  moonPhaseFor,
  scoreLabel,
  suggestedBaits,
  suggestedTechniques
} from '@/lib/fishingRules'

// Las reglas son las mismas que usa la app Android (WeatherRepository.calculateFavorabilityScore),
// así que la nota de "cuándo pescar" coincide en los dos lados.
const PLACES = [
  { name: 'Golfo de Nicoya', lat: 9.9382, lon: -84.8371 },
  { name: 'Playa Hermosa', lat: 9.9701, lon: -84.8381 },
  { name: 'Isla Chira', lat: 10.1035, lon: -85.1588 },
  { name: 'Tambor', lat: 9.9795, lon: -84.8817 }
]

const places = ref([...PLACES])
const selected = ref(PLACES[0])
const weather = ref(null)
const loading = ref(true)
const error = ref(null)
const mapContainer = ref(null)
let map = null
let marker = null

const moonPhase = computed(() => moonPhaseFor(new Date()))

const score = computed(() => {
  if (!weather.value) return null
  return favorabilityScore(weather.value, moonPhase.value)
})

const scoreLabelText = computed(() => (score.value === null ? '' : scoreLabel(score.value)))

const techniques = computed(() => {
  if (!weather.value) return []
  return suggestedTechniques(weather.value, score.value)
})

const baits = computed(() => {
  if (!weather.value) return []
  return suggestedBaits(weather.value, score.value)
})

const bestTime = computed(() => {
  if (!weather.value) return null
  return bestTimeToFish(new Date(), moonPhase.value)
})

async function loadPlaces() {
  try {
    const { data } = await axios.get('/api/v1/spots')
    const fromSpots = data.map((s) => ({
      name: s.name,
      lat: s.latitude,
      lon: s.longitude
    }))
    if (fromSpots.length) places.value = [...fromSpots, ...PLACES]
  } catch {
    // Sin spots no pasa nada: quedan los lugares por defecto.
  }
}

async function loadWeather() {
  loading.value = true
  error.value = null
  try {
    const { data } = await axios.get('/api/v1/weather', {
      params: { lat: selected.value.lat, lon: selected.value.lon }
    })
    weather.value = data
  } catch (e) {
    weather.value = null
    error.value = e.response?.data?.error || 'No se pudo obtener el clima'
  } finally {
    loading.value = false
    await nextTick()
    drawMap()
  }
}

function drawMap() {
  if (!mapContainer.value) return
  if (!map) {
    map = L.map(mapContainer.value).setView([selected.value.lat, selected.value.lon], 9)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map)
  } else {
    map.setView([selected.value.lat, selected.value.lon], 9)
  }
  if (marker) marker.remove()
  marker = L.marker([selected.value.lat, selected.value.lon]).addTo(map)
}

watch(selected, loadWeather)

onMounted(async () => {
  await loadPlaces()
  await loadWeather()
})
</script>

<template>
  <div class="weather-view">
    <div class="page-head">
      <div>
        <span class="kicker">Condiciones</span>
        <h1 class="hemmingway-title">Clima y cuándo pescar</h1>
      </div>
      <select v-model="selected" class="hemmingway-input place-select">
        <option v-for="place in places" :key="place.name" :value="place">
          {{ place.name }}
        </option>
      </select>
    </div>

    <div v-if="loading" class="loading">Consultando el clima...</div>

    <div v-else-if="error" class="empty-state">
      <h3>No hay datos del clima</h3>
      <p class="text-muted">{{ error }}</p>
      <p v-if="error.includes('OPENWEATHER_API_KEY')" class="text-muted">
        El servidor necesita la variable <code>OPENWEATHER_API_KEY</code> en su <code>.env</code>.
      </p>
    </div>

    <template v-else-if="weather">
      <div class="score-row">
        <div class="hemmingway-card score-card">
          <span class="score-number">{{ score }}</span>
          <span class="score-label">{{ scoreLabelText }}</span>
        </div>
        <div class="hemmingway-card moon-card">
          <span class="moon-phase">{{ MOON_LABELS[moonPhase] }}</span>
          <span class="text-muted">Fase aproximada</span>
        </div>
        <div v-if="bestTime" class="hemmingway-card moon-card">
          <span class="moon-phase">{{ bestTime.toLocaleString('es-CR') }}</span>
          <span class="text-muted">Mejor momento</span>
        </div>
      </div>

      <div ref="mapContainer" class="map"></div>

      <div class="cond-grid">
        <div class="hemmingway-card cond">
          <span class="cond-label">Temperatura</span>
          <strong>{{ Math.round(weather.temperature) }}°C</strong>
          <span class="text-muted">sensación {{ Math.round(weather.feelsLike) }}°C</span>
        </div>
        <div class="hemmingway-card cond">
          <span class="cond-label">Viento</span>
          <strong>{{ weather.windSpeed }} m/s</strong>
          <span class="text-muted">{{ weather.windDeg }}°</span>
        </div>
        <div class="hemmingway-card cond">
          <span class="cond-label">Nubes</span>
          <strong>{{ weather.cloudCover }}%</strong>
        </div>
        <div class="hemmingway-card cond">
          <span class="cond-label">Humedad</span>
          <strong>{{ weather.humidity }}%</strong>
        </div>
        <div class="hemmingway-card cond">
          <span class="cond-label">Presión</span>
          <strong>{{ weather.pressure }} hPa</strong>
        </div>
        <div class="hemmingway-card cond">
          <span class="cond-label">Cielo</span>
          <strong class="desc">{{ weather.description || '—' }}</strong>
        </div>
      </div>

      <div class="lists">
        <div v-if="techniques.length" class="hemmingway-card">
          <h3>Técnicas sugeridas</h3>
          <ul><li v-for="t in techniques" :key="t">{{ t }}</li></ul>
        </div>
        <div v-if="baits.length" class="hemmingway-card">
          <h3>Cebos sugeridos</h3>
          <ul><li v-for="b in baits" :key="b">{{ b }}</li></ul>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}
.place-select {
  min-width: 220px;
}
.score-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
  margin-bottom: 16px;
}
.score-card,
.moon-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.score-number {
  font-family: 'IBM Plex Mono', monospace;
  font-size: 2.2rem;
  font-weight: 700;
  color: var(--hemmingway-accent, #38bdf8);
}
.score-label {
  font-size: 0.9rem;
}
.moon-phase {
  font-size: 1.05rem;
  font-weight: 600;
}
.map {
  height: 260px;
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 16px;
  background: #0b1626;
}
.cond-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 14px;
  margin-bottom: 16px;
}
.cond {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.cond-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
.cond strong {
  font-family: 'IBM Plex Mono', monospace;
  font-size: 1.15rem;
}
.cond .desc {
  font-size: 0.95rem;
  text-transform: capitalize;
}
.lists {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 14px;
}
.lists h3 {
  margin: 0 0 10px;
}
.lists ul {
  margin: 0;
  padding-left: 18px;
}
.loading,
.empty-state {
  padding: 32px 0;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
.empty-state {
  text-align: center;
}
.empty-state code {
  font-family: 'IBM Plex Mono', monospace;
}
</style>
