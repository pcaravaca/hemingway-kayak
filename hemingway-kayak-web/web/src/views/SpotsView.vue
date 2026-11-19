<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const auth = useAuthStore()
const isAuthenticated = computed(() => auth.isAuthenticated)

const spots = ref([])
const loading = ref(true)
const error = ref(null)
const selected = ref(null)
const mapContainer = ref(null)
let map = null

const showForm = ref(false)
const saving = ref(false)
const formError = ref(null)
const form = ref({ name: '', description: '', latitude: '', longitude: '', rating: 5 })

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

async function save() {
  formError.value = null
  if (!form.value.name.trim()) {
    formError.value = 'El nombre es obligatorio'
    return
  }
  const lat = Number(form.value.latitude)
  const lon = Number(form.value.longitude)
  if (!Number.isFinite(lat) || !Number.isFinite(lon)) {
    formError.value = 'Faltan las coordenadas'
    return
  }
  saving.value = true
  try {
    await axios.post('/api/v1/spots/', {
      name: form.value.name.trim(),
      description: form.value.description || null,
      latitude: lat,
      longitude: lon,
      rating: Number(form.value.rating) || 0,
      isShared: true
    })
    showForm.value = false
    form.value = { name: '', description: '', latitude: '', longitude: '', rating: 5 }
    await load()
  } catch (e) {
    formError.value = e.response?.data?.error || 'No se pudo compartir el spot'
  } finally {
    saving.value = false
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
        <button
          v-if="isAuthenticated"
          class="hemmingway-btn hemmingway-btn-primary"
          @click="showForm = true"
        >
          Compartir spot
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">Cargando spots...</div>
    <div v-else-if="error" class="empty-state">
      <h3>No se pudieron cargar los spots</h3>
      <p class="text-muted">{{ error }}</p>
    </div>

    <div v-if="showForm" class="hemmingway-card spot-form">
      <h3>Compartir un spot</h3>
      <label class="field">
        <span>Nombre</span>
        <input v-model="form.name" class="hemmingway-input" />
      </label>
      <label class="field">
        <span>Descripción</span>
        <input v-model="form.description" class="hemmingway-input" />
      </label>
      <div class="field-row">
        <label class="field">
          <span>Latitud</span>
          <input v-model="form.latitude" type="number" step="any" class="hemmingway-input" />
        </label>
        <label class="field">
          <span>Longitud</span>
          <input v-model="form.longitude" type="number" step="any" class="hemmingway-input" />
        </label>
      </div>
      <label class="field">
        <span>Calificación (0-5)</span>
        <input v-model="form.rating" type="number" min="0" max="5" class="hemmingway-input" />
      </label>
      <div class="form-actions">
        <button class="hemmingway-btn hemmingway-btn-ghost" @click="showForm = false">Cancelar</button>
        <button class="hemmingway-btn hemmingway-btn-primary" :disabled="saving" @click="save">
          {{ saving ? 'Guardando...' : 'Compartir' }}
        </button>
      </div>
      <p v-if="formError" class="form-error">{{ formError }}</p>
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
.spot-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-bottom: 18px;
}
.spot-form h3 {
  margin: 0;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.field > span {
  font-size: 0.8rem;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.form-error {
  margin: 0;
  color: #ff8a80;
  font-size: 0.85rem;
}
</style>
