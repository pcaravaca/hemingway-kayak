<template>
  <div class="history-view">
    <div class="history-header">
      <h1 class="hemingway-title">Historial de salidas</h1>
      <div class="header-actions">
        <select v-model="filterType" class="hemingway-input filter-select">
          <option value="ALL">Todos los tipos</option>
          <option value="FISHING">Pesca</option>
          <option value="CAMPING">Campamento</option>
          <option value="TOUR">Recorrido</option>
          <option value="RACING">Carrera</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="loading">Cargando salidas...</div>

    <div v-else-if="filteredTrips.length === 0" class="no-trips">
      <div class="empty-state">
        <div class="empty-icon">🚣</div>
        <h3>Todavía no hay salidas</h3>
        <p class="text-muted">Empezá tu primera salida desde la app Android</p>
      </div>
    </div>

    <div v-else class="trips-list">
      <div
        v-for="trip in filteredTrips"
        :key="trip.id"
        class="trip-card hemingway-card"
        @click="viewTrip(trip)"
      >
        <div class="trip-header">
          <div class="trip-title-section">
            <h3 class="trip-title">{{ trip.title }}</h3>
            <span :class="['hemingway-badge', getTypeBadge(trip.type)]">
              {{ trip.type }}
            </span>
          </div>
          <span :class="['hemingway-badge', getStatusBadge(trip.status)]">
            {{ trip.status }}
          </span>
        </div>

        <div class="trip-details">
          <div class="detail-item">
            <span class="detail-icon">📅</span>
            <span>{{ formatDateRange(trip.startTime, trip.endTime) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-icon">📏</span>
            <span>{{ trip.distance ? `${trip.distance.toFixed(1)} km` : 'N/A' }}</span>
          </div>
        </div>

        <div v-if="trip.startLat && trip.startLon" class="trip-location">
          <span class="location-icon">📍</span>
          <span class="text-muted">{{ formatLocation(trip.startLat, trip.startLon) }}</span>
        </div>

        <div class="trip-footer">
          <button class="hemingway-btn hemingway-btn-secondary hemingway-btn-sm" @click.stop="viewTrip(trip)">
            View Details
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const trips = ref([])
const loading = ref(true)
const filterType = ref('ALL')

const filteredTrips = computed(() => {
  if (filterType.value === 'ALL') return trips.value
  return trips.value.filter(t => t.type === filterType.value)
})

onMounted(async () => {
  await fetchTrips()
})

async function fetchTrips() {
  loading.value = true
  try {
    // Solo las salidas del usuario logueado: /trips es el feed público (incluye
    // salidas de otros) y /trips/:id está protegido por ownership, así que
    // clickear una ajena daba 403.
    const userId = auth.user?.id
    if (!userId) {
      trips.value = []
      return
    }
    const response = await axios.get(`/api/v1/history/user/${userId}/trips`)
    trips.value = response.data.sort((a, b) => new Date(b.startTime) - new Date(a.startTime))
  } catch (error) {
    console.error('Failed to fetch trips:', error)
  } finally {
    loading.value = false
  }
}

function viewTrip(trip) {
  router.push(`/trips/${trip.id}`)
}

function getTypeBadge(type) {
  switch (type) {
    case 'FISHING': return 'hemingway-badge-success'
    case 'CAMPING': return 'hemingway-badge-warning'
    case 'TOUR': return 'hemingway-badge-info'
    case 'RACING': return 'hemingway-badge-error'
    default: return 'hemingway-badge-info'
  }
}

function getStatusBadge(status) {
  switch (status) {
    case 'ACTIVE': return 'hemingway-badge-success'
    case 'COMPLETED': return 'hemingway-badge-info'
    case 'PLANNED': return 'hemingway-badge-warning'
    case 'CANCELLED': return 'hemingway-badge-error'
    default: return 'hemingway-badge-info'
  }
}

function formatDateRange(start, end) {
  if (!start) return 'Date not set'
  const startDate = new Date(start)
  const endDate = end ? new Date(end) : null

  const dateOptions = { month: 'short', day: 'numeric', year: 'numeric' }
  const timeOptions = { hour: '2-digit', minute: '2-digit' }

  const dateStr = startDate.toLocaleDateString('en-US', dateOptions)
  const timeStr = startDate.toLocaleTimeString('en-US', timeOptions)

  if (!endDate) return `${dateStr} at ${timeStr}`

  const sameDay = startDate.toDateString() === endDate.toDateString()
  if (sameDay) {
    return `${dateStr} ${startDate.toLocaleTimeString('en-US', timeOptions)} - ${endDate.toLocaleTimeString('en-US', timeOptions)}`
  }
  return `${dateStr} - ${endDate.toLocaleDateString('en-US', dateOptions)}`
}

function formatLocation(lat, lon) {
  if (!lat || !lon) return ''
  return `${lat.toFixed(4)}, ${lon.toFixed(4)}`
}
</script>

<style scoped>
.history-view {
  max-width: 1000px;
  margin: 0 auto;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--hemmingway-space-xl);
}

.filter-select {
  width: 160px;
}

.loading, .no-trips {
  text-align: center;
  padding: var(--hemmingway-space-2xl);
}

.empty-state {
  padding: var(--hemmingway-space-2xl);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: var(--hemmingway-space-lg);
}

.empty-state h3 {
  color: var(--hemmingway-text-primary);
  margin-bottom: var(--hemmingway-space-sm);
}

.trips-list {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-lg);
}

.trip-card {
  cursor: pointer;
  transition: transform var(--hemmingway-transition-fast);
}

.trip-card:hover {
  transform: translateY(-2px);
}

.trip-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--hemmingway-space-lg);
}

.trip-title-section {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-md);
}

.trip-title {
  font-size: var(--hemmingway-font-size-lg);
  color: var(--hemmingway-text-primary);
  margin: 0;
}

.trip-details {
  display: flex;
  gap: var(--hemmingway-space-xl);
  margin-bottom: var(--hemmingway-space-md);
}

.detail-item {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-sm);
  color: var(--hemmingway-text-secondary);
  font-size: var(--hemmingway-font-size-sm);
}

.detail-icon {
  font-size: 16px;
}

.trip-location {
  display: flex;
  align-items: center;
  gap: var(--hemmingway-space-sm);
  margin-bottom: var(--hemmingway-space-lg);
}

.location-icon {
  font-size: 14px;
}

.trip-footer {
  display: flex;
  justify-content: flex-end;
}

.hemingway-btn-sm {
  padding: var(--hemmingway-space-sm) var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
}
</style>