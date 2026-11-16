<template>
  <div class="tracking-view">
    <div class="tracking-header">
      <h1 class="hemingway-title">Live Tracking</h1>
      <div class="tracking-controls">
        <button class="hemingway-btn hemingway-btn-primary" @click="connect" :disabled="wsConnected">
          {{ wsConnected ? 'Connected' : 'Connect' }}
        </button>
        <button class="hemingway-btn hemingway-btn-secondary" @click="disconnect" :disabled="!wsConnected">
          Disconnect
        </button>
      </div>
    </div>

    <div class="tracking-content">
      <div class="map-container hemingway-card">
        <div id="tracking-map" ref="mapContainer"></div>
      </div>

      <div class="trackers-panel">
        <h3 class="panel-title">Active Kayakers</h3>
        <div v-if="trackers.length === 0" class="no-trackers">
          <p class="text-muted">No kayakers sharing their location</p>
        </div>
        <div v-else class="trackers-list">
          <div
            v-for="tracker in trackers"
            :key="tracker.userId"
            class="tracker-card hemingway-card"
            @click="focusTracker(tracker)"
          >
            <div class="tracker-header">
              <span class="tracker-name">{{ tracker.displayName || tracker.username || 'Unknown' }}</span>
              <span class="hemingway-badge hemingway-badge-success">LIVE</span>
            </div>
            <div class="tracker-location">
              {{ formatLocation(tracker.latitude, tracker.longitude) }}
            </div>
            <div v-if="tracker.speed" class="tracker-speed">
              {{ tracker.speed.toFixed(1) }} km/h
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { useTrackingWebSocket } from '@/composables/useTrackingWebSocket'

const mapContainer = ref(null)
let map = null
const trackerMarkers = {}

const {
  connected: wsConnected,
  trackers,
  connect,
  disconnect
} = useTrackingWebSocket()

onMounted(() => {
  initMap()
  connect()
})

onUnmounted(() => {
  disconnect()
  if (map) {
    map.remove()
    map = null
  }
})

function initMap() {
  if (!mapContainer.value) return

  map = L.map(mapContainer.value).setView([25.7617, -80.1918], 10)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map)
}

watch(trackers, (newTrackers) => {
  updateMarkers(newTrackers)
}, { deep: true })

function updateMarkers(newTrackers) {
  newTrackers.forEach(tracker => {
    if (!map) return
    const { userId, latitude, longitude, displayName } = tracker

    if (latitude == null || longitude == null) return

    if (trackerMarkers[userId]) {
      trackerMarkers[userId].setLatLng([latitude, longitude])
    } else {
      const marker = L.marker([latitude, longitude], {
        icon: createTrackerIcon()
      }).addTo(map)
      marker.bindPopup(`<b>${displayName || 'Kayaker'}</b>`)
      trackerMarkers[userId] = marker
    }
  })
}

function createTrackerIcon() {
  return L.divIcon({
    className: 'tracker-marker',
    html: `<div style="background: #4DFFA9; width: 20px; height: 20px; border-radius: 50%; border: 3px solid #1A1B1E; box-shadow: 0 0 10px #4DFFA9;"></div>`,
    iconSize: [20, 20],
    iconAnchor: [10, 10]
  })
}

function focusTracker(tracker) {
  if (map && tracker.latitude && tracker.longitude) {
    map.setView([tracker.latitude, tracker.longitude], 14)
    if (trackerMarkers[tracker.userId]) {
      trackerMarkers[tracker.userId].openPopup()
    }
  }
}

function formatLocation(lat, lon) {
  if (!lat || !lon) return 'Unknown'
  return `${lat.toFixed(5)}, ${lon.toFixed(5)}`
}
</script>

<style scoped>
.tracking-view {
  max-width: 1400px;
  margin: 0 auto;
}

.tracking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--hemmingway-space-xl);
}

.tracking-controls {
  display: flex;
  gap: var(--hemmingway-space-md);
}

.tracking-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: var(--hemmingway-space-xl);
}

.map-container {
  height: 600px;
  padding: 0;
  overflow: hidden;
}

#tracking-map {
  width: 100%;
  height: 100%;
}

.panel-title {
  font-size: var(--hemmingway-font-size-lg);
  color: var(--hemmingway-neon);
  margin-bottom: var(--hemmingway-space-lg);
}

.trackers-list {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-md);
  max-height: 520px;
  overflow-y: auto;
}

.tracker-card {
  cursor: pointer;
  transition: transform var(--hemmingway-transition-fast);
}

.tracker-card:hover {
  transform: translateX(4px);
}

.tracker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--hemmingway-space-sm);
}

.tracker-name {
  font-weight: 600;
  color: var(--hemmingway-text-primary);
}

.tracker-location {
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-text-muted);
}

.tracker-speed {
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-turquoise);
  margin-top: var(--hemmingway-space-xs);
}

.no-trackers {
  text-align: center;
  padding: var(--hemmingway-space-xl);
}

@media (max-width: 900px) {
  .tracking-content {
    grid-template-columns: 1fr;
  }

  .map-container {
    height: 400px;
  }
}
</style>