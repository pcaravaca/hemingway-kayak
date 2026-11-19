<template>
  <div class="tournament-detail">
    <router-link to="/tournaments" class="back-btn hemingway-btn hemmingway-btn-ghost">
      ← Volver a torneos
    </router-link>

    <div v-if="loading" class="state">Cargando el torneo…</div>

    <div v-else-if="tournament" class="tournament-content">
      <div class="tournament-header">
        <div class="header-info">
          <h1>{{ tournament.name }}</h1>
          <span :class="['badge', statusClass(tournament.status)]">{{ statusLabel(tournament.status) }}</span>
        </div>
        <div class="header-meta">
          <span>{{ typeLabel(tournament.tournamenttype || tournament.type) }}</span>
          <span class="dot">•</span>
          <span>{{ formatDateRange(tournament.startTime, tournament.endTime) }}</span>
          <span v-if="tournament.locationName" class="dot">•</span>
          <span v-if="tournament.locationName">{{ tournament.locationName }}</span>
        </div>
      </div>

      <div class="tournament-grid">
        <div class="map-section">
          <div ref="mapContainer" id="tournament-map"></div>
        </div>

        <div class="rankings-section">
          <h3 class="section-title">Posiciones</h3>

          <div v-if="rankings.length === 0" class="panel-empty">
            Todavía no hay participantes.
          </div>

          <div v-else class="rankings-list">
            <div
              v-for="participant in rankings"
              :key="participant.id"
              class="ranking-item"
              @click="focusParticipant(participant)"
            >
              <div class="rank">{{ participant.rank }}</div>
              <div class="participant-info">
                <div class="participant-name">
                  {{ participant.displayName || participant.username }}
                  <span v-if="isMe(participant)" class="me">vos</span>
                </div>
                <div class="participant-stats">
                  <span>Capturas: {{ participant.catchesCount }}</span>
                  <span v-if="participant.currentLat">
                    📍 {{ formatLocation(participant.currentLat, participant.currentLon) }}
                  </span>
                </div>
              </div>
              <div class="participant-score">
                {{ participant.score != null ? participant.score.toFixed(1) : '0.0' }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Cargar captura -->
      <section v-if="myParticipation" class="panel">
        <div class="panel-head">
          <h3>Cargar una captura</h3>
          <span class="panel-note">Suma puntos al torneo</span>
        </div>
        <form class="catch-form" @submit.prevent="submitCatch">
          <label class="field">
            <span>Especie</span>
            <select v-model="catchForm.species" class="hemmingway-input" required>
              <option value="" disabled>Elegí la especie</option>
              <optgroup v-for="g in speciesGroups" :key="g.label" :label="g.label">
                <option v-for="s in g.items" :key="s.id" :value="s.commonName">
                  {{ s.commonName }}
                </option>
              </optgroup>
              <option value="__otra__">Otra (no está en la lista)</option>
            </select>
          </label>
          <label v-if="catchForm.species === '__otra__'" class="field">
            <span>Especie (a mano)</span>
            <input v-model="catchForm.speciesOther" class="hemmingway-input" placeholder="Ej.: Robalo" required />
          </label>
          <p v-if="selectedSpecies" class="field-hint">
            <span v-if="selectedSpecies.sizeCm" class="hint-strong">Talla habitual {{ selectedSpecies.sizeCm }} cm</span>
            <span class="hint-sep" v-if="selectedSpecies.sizeCm">·</span>
            <span>{{ selectedSpecies.incopesca }}</span>
          </p>
          <div class="row">
            <label class="field">
              <span>Peso (kg)</span>
              <input v-model.number="catchForm.weight" type="number" step="0.1" min="0" class="hemmingway-input" required />
            </label>
            <label class="field">
              <span>Foto (URL, opcional)</span>
              <input v-model="catchForm.photoUrl" class="hemmingway-input" placeholder="https://…" />
            </label>
          </div>
          <div class="row">
            <label class="field">
              <span>Latitud</span>
              <input v-model.number="catchForm.latitude" type="number" step="any" class="hemmingway-input" />
            </label>
            <label class="field">
              <span>Longitud</span>
              <input v-model.number="catchForm.longitude" type="number" step="any" class="hemmingway-input" />
            </label>
          </div>
          <p v-if="catchMessage" class="message" :class="{ ok: catchSaved }">{{ catchMessage }}</p>
          <div class="actions">
            <button type="submit" class="hemmingway-btn hemmingway-btn-primary" :disabled="savingCatch">
              {{ savingCatch ? 'Guardando…' : 'Cargar captura' }}
            </button>
            <button
              type="button"
              class="hemmingway-btn hemmingway-btn-secondary"
              @click="useTournamentLocation"
            >
              Usar la ubicación del torneo
            </button>
          </div>
        </form>
      </section>

      <section v-else class="panel">
        <div class="panel-empty">
          No estás anotado en este torneo, así que no podés cargar capturas.
        </div>
      </section>

      <!-- Listado de capturas -->
      <section class="panel">
        <div class="panel-head">
          <h3>Capturas</h3>
          <span class="panel-note">{{ catches.length }} registradas</span>
        </div>
        <div v-if="catches.length === 0" class="panel-empty">
          Todavía no hay capturas cargadas.
        </div>
        <ul v-else class="catches">
          <li v-for="c in catches" :key="c.id">
            <div class="catch-main">
              <b>{{ c.species }}</b>
              <span class="catch-by">{{ c.username }}</span>
              <span v-if="speciesScientific(c.species)" class="catch-sci">{{ speciesScientific(c.species) }}</span>
            </div>
            <div class="catch-meta">
              <span class="weight">{{ c.weight.toFixed(1) }} kg</span>
              <span v-if="c.verified" class="badge is-ok">VERIFICADA</span>
              <span class="time">{{ formatDate(c.createdAt) }}</span>
            </div>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { useTrackingWebSocket } from '@/composables/useTrackingWebSocket'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const auth = useAuthStore()

const tournament = ref(null)
const rankings = ref([])
const catches = ref([])
const speciesList = ref([])
const loading = ref(true)
const mapContainer = ref(null)
let map = null
const participantMarkers = {}

const savingCatch = ref(false)
const catchSaved = ref(false)
const catchMessage = ref('')

const catchForm = reactive({
  species: '',
  speciesOther: '',
  weight: null,
  latitude: null,
  longitude: null,
  photoUrl: ''
})

// Agrupa la DB de peces por hábitat para el selector del formulario de captura.
const speciesGroups = computed(() => {
  const labels = {
    rio: 'Río (agua dulce)',
    estuario: 'Estero / manglar (salobre)',
    manglar: 'Manglar',
    mar: 'Mar'
  }
  return ['rio', 'estuario', 'manglar', 'mar']
    .map(h => ({ label: labels[h], items: speciesList.value.filter(s => s.habitat === h) }))
    .filter(g => g.items.length)
})

function speciesScientific(commonName) {
  const s = speciesList.value.find(x => x.commonName === commonName)
  return s?.scientific || ''
}

// Especie elegida en el selector: de donde sale la talla legal (INCOPESCA) y la
// talla habitual que se muestran bajo el campo.
const selectedSpecies = computed(
  () => speciesList.value.find(x => x.commonName === catchForm.species) || null
)

const { connect, subscribeTournament, unsubscribeTournament, tournamentUpdates } = useTrackingWebSocket()

const myParticipation = computed(() => {
  const me = auth.user?.id
  if (!me) return null
  return rankings.value.find(p => p.userId === me) || null
})

onMounted(async () => {
  await fetchTournament()
  await nextTick()
  initMap()
  connect()
})

onUnmounted(() => {
  unsubscribeTournament(route.params.id)
  if (map) {
    map.remove()
    map = null
  }
})

async function fetchTournament() {
  loading.value = true
  try {
    const [tournamentRes, rankingRes, catchesRes, speciesRes] = await Promise.allSettled([
      axios.get(`/api/v1/tournaments/${route.params.id}`),
      axios.get(`/api/v1/tournaments/${route.params.id}/ranking`),
      axios.get(`/api/v1/tournaments/${route.params.id}/catches`),
      axios.get('/api/v1/species')
    ])
    if (tournamentRes.status === 'fulfilled') {
      tournament.value = tournamentRes.value.data
      if (tournament.value?.status === 'LIVE') subscribeTournament(route.params.id)
    }
    if (rankingRes.status === 'fulfilled') rankings.value = rankingRes.value.data || []
    if (catchesRes.status === 'fulfilled') catches.value = catchesRes.value.data || []
    if (speciesRes.status === 'fulfilled') speciesList.value = speciesRes.value.data || []
  } catch (error) {
    console.error('No se pudo cargar el torneo:', error)
  } finally {
    loading.value = false
  }
}

async function submitCatch() {
  const participation = myParticipation.value
  if (!participation) {
    catchMessage.value = 'No estás anotado en este torneo.'
    catchSaved.value = false
    return
  }

  savingCatch.value = true
  catchMessage.value = ''
  // Si eligió "Otra", el valor es el que tipeó a mano; si no, el nombre común
  // de la DB (que es lo que se guarda en tournament_catches.species).
  const species = catchForm.species === '__otra__'
    ? (catchForm.speciesOther || '').trim()
    : catchForm.species
  if (!species) {
    catchSaved.value = false
    catchMessage.value = 'Elegí o escribí la especie.'
    savingCatch.value = false
    return
  }
  try {
    await axios.post(`/api/v1/tournaments/${route.params.id}/catches`, {
      participantId: participation.id,
      species,
      weight: Number(catchForm.weight),
      latitude: catchForm.latitude,
      longitude: catchForm.longitude,
      photoUrl: catchForm.photoUrl || null
    })
    catchSaved.value = true
    catchMessage.value = 'Captura cargada.'
    catchForm.species = ''
    catchForm.speciesOther = ''
    catchForm.weight = null
    catchForm.photoUrl = ''
    await fetchTournament()
    // El mapa se volvió a crear: hay que repintar los marcadores.
    await nextTick()
    initMap()
  } catch (e) {
    catchSaved.value = false
    catchMessage.value = e.response?.data?.error || 'No se pudo cargar la captura.'
  } finally {
    savingCatch.value = false
  }
}

function useTournamentLocation() {
  if (tournament.value?.startLat != null) {
    catchForm.latitude = tournament.value.startLat
    catchForm.longitude = tournament.value.startLon
  }
}

watch(tournamentUpdates, (updates) => {
  if (updates?.rankings) {
    rankings.value = updates.rankings
    updateMarkers()
  }
})

function isMe(participant) {
  return auth.user?.id && participant.userId === auth.user.id
}

function initMap() {
  if (!mapContainer.value) return
  if (map) {
    map.remove()
    map = null
  }

  const t = tournament.value
  const center = t?.startLat != null
    ? [t.startLat, t.startLon]
    : [9.9382, -84.8371]

  map = L.map(mapContainer.value).setView(center, 11)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map)

  if (t?.startLat != null) {
    L.circleMarker(center, {
      radius: 8, fillColor: '#38BDF8', color: '#060D1A', weight: 2, fillOpacity: 1
    }).addTo(map).bindPopup(t.locationName || 'Torneo')
  }

  updateMarkers()
}

function updateMarkers() {
  rankings.value.forEach(participant => {
    if (!map || !participant.currentLat || !participant.currentLon) return

    const markerId = participant.id
    if (participantMarkers[markerId]) {
      participantMarkers[markerId].setLatLng([participant.currentLat, participant.currentLon])
    } else {
      const marker = L.marker([participant.currentLat, participant.currentLon], {
        icon: createParticipantIcon(participant.rank)
      }).addTo(map)
      marker.bindPopup(
        `<b>${participant.displayName || participant.username}</b><br>Score: ${participant.score ?? 0}`
      )
      participantMarkers[markerId] = marker
    }
  })
}

function createParticipantIcon(rank) {
  const colors = ['#4DFFA9', '#40E0D0', '#FFD700', '#FF6B6B']
  const color = rank <= 4 ? colors[rank - 1] : '#707070'
  return L.divIcon({
    className: 'participant-marker',
    html: `<div style="background: ${color}; width: 24px; height: 24px; border-radius: 50%; border: 3px solid #1A1B1E; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: bold; color: #1A1B1E;">${rank}</div>`,
    iconSize: [24, 24],
    iconAnchor: [12, 12]
  })
}

function focusParticipant(participant) {
  if (map && participant.currentLat && participant.currentLon) {
    map.setView([participant.currentLat, participant.currentLon], 14)
    if (participantMarkers[participant.id]) participantMarkers[participant.id].openPopup()
  }
}

function statusLabel(status) {
  return {
    PLANNED: 'PLANIFICADO',
    UPCOMING: 'PRÓXIMO',
    REGISTRATION: 'INSCRIPCIÓN',
    LIVE: 'EN VIVO',
    FINISHED: 'FINALIZADO',
    CANCELLED: 'CANCELADO'
  }[status] || (status || '')
}

function statusClass(status) {
  if (status === 'LIVE') return 'is-live'
  if (status === 'FINISHED') return 'is-done'
  return 'is-next'
}

function typeLabel(type) {
  return { FISHING: 'Pesca', RACING: 'Carrera', TOURING: 'Travesía' }[type] || type || 'Torneo'
}

function formatDateRange(start, end) {
  if (!start) return 'Fecha a definir'
  const startDate = new Date(start)
  if (Number.isNaN(startDate.getTime())) return 'Fecha a definir'
  const fmt = { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' }
  const startStr = startDate.toLocaleDateString('es-CR', fmt)
  if (!end) return startStr
  const endDate = new Date(end)
  if (Number.isNaN(endDate.getTime())) return startStr
  return `${startStr} — ${endDate.toLocaleDateString('es-CR', { hour: '2-digit', minute: '2-digit' })}`
}

function formatDate(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('es-CR', { day: '2-digit', month: 'short' })
}

function formatLocation(lat, lon) {
  if (!lat || !lon) return ''
  return `${lat.toFixed(4)}, ${lon.toFixed(4)}`
}
</script>

<style scoped>
.tournament-detail { padding: 40px 0 24px; }

.state { padding: 64px 0; color: var(--hemmingway-text-muted); font-size: 15px; }

.tournament-content {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-xl);
}

.tournament-header {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-md);
}

.header-info {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-lg);
  flex-wrap: wrap;
}

.header-info h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 38px;
  letter-spacing: -.02em;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-md);
  color: var(--hemmingway-text-secondary);
  font-size: 14px;
  flex-wrap: wrap;
}

.dot { color: var(--hemmingway-text-muted); }

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
.badge.is-ok { border-color: var(--hemmingway-accent-success); color: var(--hemmingway-accent-success); }

.tournament-grid {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: var(--hemmingway-space-xl);
}

.map-section {
  height: 520px;
  border-radius: var(--hemmingway-radius-xl);
  overflow: hidden;
  border: 1px solid var(--hemmingway-line);
}

#tournament-map { width: 100%; height: 100%; }

.section-title {
  font-family: var(--hemmingway-font-display);
  font-weight: 500;
  font-size: 22px;
  margin-bottom: var(--hemmingway-space-lg);
}

.rankings-list {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-md);
  max-height: 480px;
  overflow-y: auto;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-lg);
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 14px 16px;
  cursor: pointer;
  transition: border-color var(--hemmingway-transition-fast), transform var(--hemmingway-transition-fast);
}

.ranking-item:hover {
  border-color: var(--hemmingway-blue);
  transform: translateX(4px);
}

.rank {
  width: 32px;
  height: 32px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--hemmingway-blue);
  color: var(--hemmingway-dark);
  border-radius: 50%;
  font-weight: 700;
  font-size: 13px;
}

.participant-info { flex: 1; min-width: 0; }

.participant-name { font-weight: 600; color: var(--hemmingway-text-primary); }

.me {
  font-family: var(--hemmingway-font-mono);
  font-size: 9.5px;
  letter-spacing: .08em;
  text-transform: uppercase;
  color: var(--hemmingway-blue);
  border: 1px solid var(--hemmingway-blue);
  border-radius: var(--hemmingway-radius-full);
  padding: 2px 6px;
  margin-left: 8px;
  vertical-align: middle;
}

.participant-stats {
  display: flex;
  gap: var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-text-muted);
  margin-top: 4px;
}

.participant-score {
  font-family: var(--hemmingway-font-display);
  font-size: 22px;
  font-weight: 500;
  color: var(--hemmingway-gold);
}

.panel {
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 22px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
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

.panel-empty { color: var(--hemmingway-text-muted); font-size: 14px; }

.catch-form { display: flex; flex-direction: column; gap: 16px; }

.row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

.field { display: flex; flex-direction: column; gap: 8px; }

.field span {
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
}

/* Talla habitual + talla legal (INCOPESCA) de la especie elegida. Va fuera del
   .field para que no quede como label del selector. */
.field-hint {
  margin: -8px 0 0;
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  line-height: 1.5;
  color: var(--hemmingway-text-muted);
}

.field-hint .hint-strong { color: var(--hemmingway-accent, #38BDF8); }
.field-hint .hint-sep { margin: 0 6px; opacity: .5; }

.message { font-size: 14px; color: var(--hemmingway-accent-error); margin: 0; }
.message.ok { color: var(--hemmingway-accent-success); }

.actions { display: flex; gap: 12px; flex-wrap: wrap; }

.catches { list-style: none; margin: 0; padding: 0; }

.catches li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 13px 0;
  border-bottom: 1px solid var(--hemmingway-line);
}

.catches li:last-child { border-bottom: none; }

.catch-main b { display: block; font-size: 15px; }
.catch-by {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  color: var(--hemmingway-text-muted);
}

.catch-sci {
  display: block;
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  font-style: italic;
  color: var(--hemmingway-text-muted);
  opacity: 0.8;
}

.catch-meta {
  display: flex;
  align-items: center;
  gap: 14px;
  flex: none;
}

.catch-meta .weight {
  font-family: var(--hemmingway-font-mono);
  color: var(--hemmingway-gold);
  font-size: 14px;
}

.catch-meta .time {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  color: var(--hemmingway-text-muted);
}

@media (max-width: 1000px) {
  .tournament-grid { grid-template-columns: 1fr; }
  .map-section { height: 380px; }
  .header-info h1 { font-size: 30px; }
}

@media (max-width: 620px) {
  .tournament-detail { padding: 28px 0 24px; }
  .row { grid-template-columns: 1fr; }
  .catches li { flex-direction: column; align-items: flex-start; gap: 6px; }
}
</style>
