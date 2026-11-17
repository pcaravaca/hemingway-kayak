<template>
  <div class="trip">
    <router-link to="/history" class="back">← Volver al historial</router-link>

    <div v-if="loading" class="state">Cargando la salida…</div>

    <div v-else-if="error" class="state error">{{ error }}</div>

    <template v-else-if="trip">
      <header class="trip-head">
        <div>
          <span class="kicker"><i></i>Salida</span>
          <h1>{{ trip.title || 'Salida sin título' }}</h1>
          <p v-if="trip.description" class="lede">{{ trip.description }}</p>
          <div class="tags">
            <span class="badge" :class="statusClass">{{ statusLabel }}</span>
            <span class="badge">{{ typeLabel }}</span>
          </div>
        </div>
      </header>

      <div class="kpis">
        <div class="kpi">
          <span class="kpi-label">Distancia</span>
          <b class="kpi-value">{{ distanceLabel }}</b>
        </div>
        <div class="kpi">
          <span class="kpi-label">Inicio</span>
          <b class="kpi-value sm">{{ formatDate(trip.startTime) }}</b>
        </div>
        <div class="kpi">
          <span class="kpi-label">Fin</span>
          <b class="kpi-value sm">{{ trip.endTime && trip.endTime.startsWith('0001') ? '—' : formatDate(trip.endTime) }}</b>
        </div>
        <div class="kpi">
          <span class="kpi-label">Puntos</span>
          <b class="kpi-value">{{ points.length }}</b>
        </div>
      </div>

      <section class="panel">
        <div class="panel-head">
          <h3>Recorrido</h3>
          <span v-if="trip.startLat != null" class="panel-note">
            {{ trip.startLat.toFixed(4) }}, {{ trip.startLon.toFixed(4) }}
          </span>
        </div>
        <div ref="mapContainer" class="map"></div>
        <div v-if="points.length < 2" class="panel-empty">
          Esta salida tiene {{ points.length }} punto{{ points.length === 1 ? '' : 's' }} registrado;
          no alcanza para dibujar un recorrido.
        </div>
      </section>

      <section class="panel">
        <div class="panel-head">
          <h3>Puntos del track</h3>
          <span class="panel-note">{{ points.length }} registros</span>
        </div>
        <div v-if="!points.length" class="panel-empty">Sin puntos cargados.</div>
        <ul v-else class="points">
          <li v-for="(p, i) in points" :key="p.id">
            <span class="idx">{{ String(i + 1).padStart(2, '0') }}</span>
            <span class="coords">{{ p.latitude.toFixed(4) }}, {{ p.longitude.toFixed(4) }}</span>
            <span class="meta">
              <b v-if="p.speed != null">{{ p.speed.toFixed(1) }} kn</b>
              <em v-if="p.heading != null">· {{ Math.round(p.heading) }}°</em>
            </span>
            <span class="time">{{ formatTime(p.timestamp) }}</span>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const route = useRoute()
const loading = ref(true)
const error = ref(null)
const trip = ref(null)
const points = ref([])
const mapContainer = ref(null)
let map = null

const distanceLabel = computed(() => {
  const d = Number(trip.value?.distance) || 0
  if (!d) return '—'
  return d >= 10 ? `${d.toFixed(0)} km` : `${d.toFixed(1)} km`
})

const statusLabel = computed(() => ({
  LIVE: 'EN VIVO',
  COMPLETED: 'COMPLETADA',
  PLANNED: 'PLANIFICADA',
  CANCELLED: 'CANCELADA'
}[trip.value?.status] || trip.value?.status || ''))

const statusClass = computed(() => {
  if (trip.value?.status === 'LIVE') return 'is-live'
  if (trip.value?.status === 'COMPLETED') return 'is-done'
  return 'is-next'
})

const typeLabel = computed(() => ({
  FISHING: 'Pesca',
  TOURING: 'Travesía',
  RACING: 'Carrera',
  CAMPING: 'Campamento'
}[trip.value?.type] || trip.value?.type || 'Salida'))

onMounted(load)
onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})

async function load() {
  loading.value = true
  error.value = null
  const id = route.params.id
  try {
    const [tripRes, pointsRes] = await Promise.allSettled([
      axios.get(`/api/v1/trips/${id}`),
      axios.get(`/api/v1/trips/${id}/points`)
    ])

    if (tripRes.status === 'rejected') {
      error.value = tripRes.reason?.response?.data?.error || 'No se pudo cargar la salida.'
    } else {
      trip.value = tripRes.value.data
      points.value = pointsRes.status === 'fulfilled' ? (pointsRes.value.data || []) : []
    }
  } finally {
    loading.value = false
  }

  // Recién ahora se monta la rama del viaje: si dibujamos antes, el div del
  // mapa todavía no existe (sigue el v-if de "cargando") y Leaflet no inicializa.
  await nextTick()
  drawMap()
}

function drawMap() {
  if (!mapContainer.value || !trip.value) return
  if (map) {
    map.remove()
    map = null
  }

  const track = points.value
    .filter(p => p.latitude != null && p.longitude != null)
    .map(p => [p.latitude, p.longitude])

  const start = [trip.value.startLat, trip.value.startLon]
  const hasStart = trip.value.startLat != null && trip.value.startLon != null
  const center = track.length ? track[0] : (hasStart ? start : [9.9382, -84.8371])

  map = L.map(mapContainer.value).setView(center, track.length ? 13 : 11)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map)

  if (track.length > 1) {
    L.polyline(track, { color: '#38BDF8', weight: 3, opacity: .85 }).addTo(map)
    map.fitBounds(L.latLngBounds(track).pad(0.15))
  }

  if (track.length) {
    L.circleMarker(track[0], dotStyle('#4DFFA9')).addTo(map).bindPopup('Inicio')
    if (track.length > 1) {
      L.circleMarker(track[track.length - 1], dotStyle('#FFD700')).addTo(map).bindPopup('Fin')
    }
  } else if (hasStart) {
    L.circleMarker(start, dotStyle('#38BDF8')).addTo(map).bindPopup('Inicio')
  }
}

function dotStyle(color) {
  return {
    radius: 6,
    fillColor: color,
    color: '#060D1A',
    weight: 2,
    fillOpacity: 1
  }
}

function formatDate(value) {
  if (!value || value.startsWith('0001')) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleDateString('es-CR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatTime(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleTimeString('es-CR', { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.trip { padding: 40px 0 24px; }

.back {
  display: inline-block;
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  letter-spacing: .06em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
  text-decoration: none;
  margin-bottom: 22px;
}

.back:hover { color: var(--hemmingway-blue); }

.trip-head {
  padding-bottom: 28px;
  border-bottom: 1px solid var(--hemmingway-line);
}

.trip-head h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 40px;
  letter-spacing: -.02em;
  line-height: 1.08;
  margin: 12px 0 14px;
}

.lede {
  color: var(--hemmingway-text-secondary);
  font-size: 16px;
  line-height: 1.6;
  max-width: 62ch;
}

.tags { display: flex; gap: 8px; margin-top: 18px; flex-wrap: wrap; }

.badge {
  font-family: var(--hemmingway-font-mono);
  font-size: 10.5px;
  letter-spacing: .1em;
  padding: 6px 11px;
  border-radius: var(--hemmingway-radius-full);
  border: 1px solid var(--hemmingway-line-2);
  color: var(--hemmingway-text-secondary);
}

.badge.is-live { border-color: var(--hemmingway-neon); color: var(--hemmingway-neon); }
.badge.is-done { color: var(--hemmingway-text-muted); }
.badge.is-next { color: var(--hemmingway-blue); border-color: var(--hemmingway-blue); }

.state { padding: 64px 0; color: var(--hemmingway-text-muted); font-size: 15px; }
.state.error { color: var(--hemmingway-accent-error); }

.kpis {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 28px 0;
}

.kpi {
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 20px;
}

.kpi-label {
  display: block;
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
  margin-bottom: 12px;
}

.kpi-value {
  font-family: var(--hemmingway-font-display);
  font-weight: 500;
  font-size: 34px;
  line-height: 1;
}

.kpi-value.sm { font-size: 20px; }

.panel {
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 22px;
  margin-bottom: 16px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.panel-head h3 {
  font-family: var(--hemmingway-font-display);
  font-weight: 500;
  font-size: 19px;
}

.panel-note {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  color: var(--hemmingway-text-muted);
}

.panel-empty { color: var(--hemmingway-text-muted); font-size: 14px; padding: 8px 0; }

.map {
  height: 340px;
  border-radius: var(--hemmingway-radius-md);
  border: 1px solid var(--hemmingway-line-2);
  overflow: hidden;
  background: var(--hemmingway-bg-tertiary);
}

.points { list-style: none; margin: 0; padding: 0; }

.points li {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 11px 0;
  border-bottom: 1px solid var(--hemmingway-line);
  font-family: var(--hemmingway-font-mono);
  font-size: 12.5px;
}

.points li:last-child { border-bottom: none; }

.points .idx { color: var(--hemmingway-text-muted); flex: none; }
.points .coords { color: var(--hemmingway-text-primary); flex: 1; }
.points .meta { color: var(--hemmingway-gold); flex: none; }
.points .meta em { font-style: normal; color: var(--hemmingway-text-secondary); }
.points .time { color: var(--hemmingway-text-muted); flex: none; }

@media (max-width: 900px) {
  .trip { padding: 28px 0 24px; }
  .trip-head h1 { font-size: 30px; }
  .kpis { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 560px) {
  .kpis { grid-template-columns: 1fr; }
  .points li { flex-wrap: wrap; gap: 8px; }
}
</style>
