<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const spots = ref([])
const loading = ref(true)
const error = ref(null)
const selected = ref(null)
const mapContainer = ref(null)
let map = null

const FAVORITES = [
  { name: 'Golfo de Nicoya', lat: 9.9382, lon: -84.8371 },
  { name: 'Playa Hermosa', lat: 9.9701, lon: -84.8381 },
  { name: 'Isla Chira', lat: 10.1035, lon: -85.1588 },
  { name: 'Tambor', lat: 9.9795, lon: -84.8817 }
]

async function load() {
  loading.value = true
  try {
    const { data } = await axios.get('/api/v1/spots')
    spots.value = data
  } catch (e) {
    error.value = e.response?.data?.error || 'No se pudieron cargar los spots'
  } finally {
    loading.value = false
    await nextTick()
    initMap()
  }
}

function initMap() {
  if (!mapContainer.value || map) return
  map = L.map(mapContainer.value).setView([9.9382, -84.8371], 9)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map)

  for (const spot of spots.value) {
    if (!spot.latitude || !spot.longitude) continue
    L.marker([spot.latitude, spot.longitude])
      .addTo(map)
      .bindPopup(
        `<strong>${spot.name}</strong><br/>${spot.description || ''}<br/>` +
          `<small>${spot.displayName || spot.username || ''} · ${spot.rating || 0}/5</small>`
      )
  }
}

function flyTo(spot) {
  selected.value = spot
  if (map && spot.latitude && spot.longitude) {
    map.setView([spot.latitude, spot.longitude], 12)
  }
}

onMounted(load)

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<template>
  <div class="spots-view">
    <div class="page-head">
      <div>
        <span class="kicker">Dónde pescar</span>
        <h1 class="hemmingway-title">Spots compartidos</h1>
        <p class="text-muted">Los puntos que la comunidad marcó como buenos.</p>
      </div>
      <div class="quick">
        <button
          v-for="place in FAVORITES"
          :key="place.name"
          class="hemmingway-btn hemmingway-btn-ghost"
          @click="flyTo({ ...place, name: place.name })"
        >
          {{ place.name }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">Cargando spots...</div>
    <div v-else-if="error" class="empty-state">
      <h3>No se pudieron cargar los spots</h3>
      <p class="text-muted">{{ error }}</p>
    </div>

    <template v-else>
      <div ref="mapContainer" class="map"></div>

      <div v-if="spots.length === 0" class="empty-state">
        <h3>Todavía no hay spots compartidos</h3>
        <p class="text-muted">Compartilos desde la app Android.</p>
      </div>

      <div v-else class="spot-list">
        <button
          v-for="spot in spots"
          :key="spot.id"
          :class="['hemmingway-card', 'spot-row', { active: selected && selected.id === spot.id }]"
          @click="flyTo(spot)"
        >
          <div class="spot-main">
            <h3>{{ spot.name }}</h3>
            <p class="text-muted">{{ spot.description || 'Sin descripción' }}</p>
          </div>
          <div class="spot-meta">
            <span class="coords">{{ spot.latitude.toFixed(3) }}, {{ spot.longitude.toFixed(3) }}</span>
            <span class="hemmingway-badge">{{ spot.rating }}/5</span>
            <span class="fish">{{ spot.fishCaught }} capturas</span>
          </div>
        </button>
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
.quick {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.map {
  height: 340px;
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
  background: #0b1626;
}
.spot-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.spot-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  text-align: left;
  cursor: pointer;
  width: 100%;
}
.spot-row.active {
  border-color: var(--hemmingway-accent, #38bdf8);
}
.spot-main h3 {
  margin: 0 0 4px;
}
.spot-main p {
  margin: 0;
  font-size: 0.85rem;
}
.spot-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  font-size: 0.8rem;
}
.loading,
.empty-state {
  padding: 32px 0;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
.empty-state {
  text-align: center;
}
</style>
