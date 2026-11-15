<template>
  <div class="tournament-detail">
    <button class="back-btn hemingway-btn hemingway-btn-ghost" @click="goBack">
      ← Back to Tournaments
    </button>

    <div v-if="loading" class="loading">Loading tournament...</div>

    <div v-else-if="tournament" class="tournament-content">
      <div class="tournament-header">
        <div class="header-info">
          <h1 class="hemingway-title">{{ tournament.name }}</h1>
          <span :class="['hemingway-badge', getStatusBadge(tournament.status)]">
            {{ tournament.status }}
          </span>
        </div>
        <div class="header-meta">
          <span class="meta-item">{{ tournament.type }}</span>
          <span class="meta-separator">•</span>
          <span class="meta-item">{{ formatDateRange(tournament.startTime, tournament.endTime) }}</span>
        </div>
      </div>

      <div class="tournament-grid">
        <div class="map-section">
          <div id="tournament-map" ref="mapContainer"></div>
        </div>

        <div class="rankings-section">
          <h2 class="section-title">Rankings</h2>

          <div v-if="rankings.length === 0" class="no-data">
            <p class="text-muted">No participants yet</p>
          </div>

          <div v-else class="rankings-list">
            <div
              v-for="(participant, index) in rankings"
              :key="participant.id"
              class="ranking-item hemingway-card"
              @click="focusParticipant(participant)"
            >
              <div class="rank">{{ index + 1 }}</div>
              <div class="participant-info">
                <div class="participant-name">{{ participant.displayName || participant.username }}</div>
                <div class="participant-stats">
                  <span>Catches: {{ participant.catchesCount }}</span>
                  <span v-if="participant.currentLat">📍 {{ formatLocation(participant.currentLat, participant.currentLon) }}</span>
                </div>
              </div>
              <div class="participant-score">
                {{ participant.score.toFixed(1) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { useTrackingWebSocket } from '@/composables/useTrackingWebSocket'

const route = useRoute()
const router = useRouter()

const tournament = ref(null)
const rankings = ref([])
const loading = ref(true)
const mapContainer = ref(null)
let map = null
const participantMarkers = {}

const { subscribeTournament, unsubscribeTournament, tournamentUpdates } = useTrackingWebSocket()

onMounted(async () => {
  await fetchTournament()
  initMap()
})

onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})

async function fetchTournament() {
  loading.value = true
  try {
    const [tournamentRes, rankingRes] = await Promise.all([
      axios.get(`/api/v1/tournaments/${route.params.id}`),
      axios.get(`/api/v1/tournaments/${route.params.id}/ranking`)
    ])
    tournament.value = tournamentRes.data
    rankings.value = rankingRes.data

    if (tournament.value?.status === 'LIVE') {
      subscribeTournament(route.params.id)
    }
  } catch (error) {
    console.error('Failed to fetch tournament:', error)
  } finally {
    loading.value = false
  }
}

watch(tournamentUpdates, (updates) => {
  if (updates?.rankings) {
    rankings.value = updates.rankings
    updateMarkers()
  }
})

function initMap() {
  if (!mapContainer.value) return

  map = L.map(mapContainer.value).setView([25.7617, -80.1918], 10)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map)
}

function updateMarkers() {
  rankings.value.forEach(participant => {
    if (!map || !participant.currentLat || !participant.currentLon) return

    const markerId = participant.id

    if (participantMarkers[markerId]) {
      participantMarkers[markerId].setLatLng([participant.currentLat, participant.currentLon])
    } else {
      const marker = L.marker([participant.currentLat, participant.currentLon], {
        icon: createParticipantIcon(participant.rank || rankings.value.indexOf(participant) + 1)
      }).addTo(map)
      marker.bindPopup(`<b>${participant.displayName || participant.username}</b><br>Score: ${participant.score}`)
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
    if (participantMarkers[participant.id]) {
      participantMarkers[participant.id].openPopup()
    }
  }
}

function goBack() {
  router.push('/tournaments')
}

function getStatusBadge(status) {
  switch (status) {
    case 'LIVE': return 'hemingway-badge-success'
    case 'PLANNED': return 'hemingway-badge-info'
    case 'FINISHED': return 'hemingway-badge-warning'
    default: return 'hemingway-badge-info'
  }
}

function formatDateRange(start, end) {
  if (!start) return 'TBD'
  const startDate = new Date(start)
  const endDate = end ? new Date(end) : null

  const options = { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }
  const startStr = startDate.toLocaleDateString('en-US', options)

  if (!endDate) return startStr
  return `${startStr} - ${endDate.toLocaleDateString('en-US', { hour: '2-digit', minute: '2-digit' })}`
}

function formatLocation(lat, lon) {
  if (!lat || !lon) return ''
  return `${lat.toFixed(4)}, ${lon.toFixed(4)}`
}
</script>

<style scoped>
.tournament-detail {
  max-width: 1400px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: var(--hemmingway-space-lg);
}

.loading {
  text-align: center;
  padding: var(--hemmingway-space-2xl);
}

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
}

.header-meta {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-md);
  color: var(--hemmingway-text-secondary);
}

.meta-separator {
  color: var(--hemmingway-text-muted);
}

.tournament-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: var(--hemmingway-space-xl);
}

.map-section {
  height: 600px;
  border-radius: var(--hemmingway-radius-xl);
  overflow: hidden;
}

#tournament-map {
  width: 100%;
  height: 100%;
}

.section-title {
  font-size: var(--hemmingway-font-size-xl);
  color: var(--hemmingway-neon);
  margin-bottom: var(--hemmingway-space-lg);
}

.rankings-list {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-md);
  max-height: 560px;
  overflow-y: auto;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-lg);
  cursor: pointer;
  transition: transform var(--hemmingway-transition-fast);
}

.ranking-item:hover {
  transform: translateX(4px);
}

.rank {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--hemmingway-neon);
  color: var(--hemmingway-dark);
  border-radius: 50%;
  font-weight: 700;
}

.participant-info {
  flex: 1;
}

.participant-name {
  font-weight: 600;
  color: var(--hemmingway-text-primary);
}

.participant-stats {
  display: flex;
  gap: var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-text-muted);
}

.participant-score {
  font-size: var(--hemmingway-font-size-xl);
  font-weight: 700;
  color: var(--hemmingway-neon);
}

.no-data {
  text-align: center;
  padding: var(--hemmingway-space-xl);
}

@media (max-width: 1000px) {
  .tournament-grid {
    grid-template-columns: 1fr;
  }

  .map-section {
    height: 400px;
  }
}
</style>