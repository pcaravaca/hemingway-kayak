<template>
  <div class="tournaments-view">
    <div class="tournaments-header">
      <h1 class="hemingway-title">Torneos</h1>
      <router-link to="/tournaments/create" class="hemingway-btn hemingway-btn-primary">
        Create Tournament
      </router-link>
    </div>

    <div class="filters">
      <button
        v-for="filter in filters"
        :key="filter.value"
        :class="['filter-btn', { active: activeFilter === filter.value }]"
        @click="activeFilter = filter.value"
      >
        {{ filter.label }}
      </button>
    </div>

    <div v-if="loading" class="loading">Cargando torneos...</div>

    <div v-else-if="filteredTournaments.length === 0" class="no-tournaments">
      <p class="text-muted">No hay torneos</p>
    </div>

    <div v-else class="tournaments-grid">
      <div
        v-for="tournament in filteredTournaments"
        :key="tournament.id"
        class="tournament-card hemingway-card"
        @click="goToDetail(tournament.id)"
      >
        <div class="tournament-header">
          <h3 class="tournament-name">{{ tournament.name }}</h3>
          <span :class="['hemingway-badge', getStatusBadge(tournament.status)]">
            {{ tournament.status }}
          </span>
        </div>

        <div class="tournament-info">
          <div class="info-row">
            <span class="info-label">Tipo:</span>
            <span class="info-value">{{ tournament.type }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Empieza:</span>
            <span class="info-value">{{ formatDate(tournament.startTime) }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Participantes:</span>
            <span class="info-value">{{ tournament.participantsCount || 0 }}</span>
          </div>
        </div>

        <div class="tournament-actions">
          <button
            v-if="tournament.status === 'LIVE'"
            class="hemingway-btn hemingway-btn-primary hemingway-btn-sm"
            @click.stop="watchLive(tournament.id)"
          >
            Watch Live
          </button>
          <button
            v-else-if="tournament.status === 'PLANNED'"
            class="hemingway-btn hemingway-btn-secondary hemingway-btn-sm"
            @click.stop="joinTournament(tournament.id)"
          >
            Join
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

const router = useRouter()
const tournaments = ref([])
const loading = ref(true)
const activeFilter = ref('ALL')

const filters = [
  { label: 'All', value: 'ALL' },
  { label: 'Live', value: 'LIVE' },
  { label: 'Planned', value: 'PLANNED' },
  { label: 'Finished', value: 'FINISHED' }
]

const filteredTournaments = computed(() => {
  if (activeFilter.value === 'ALL') return tournaments.value
  return tournaments.value.filter(t => t.status === activeFilter.value)
})

onMounted(async () => {
  await fetchTournaments()
})

async function fetchTournaments() {
  loading.value = true
  try {
    const response = await axios.get('/api/v1/tournaments')
    tournaments.value = response.data
  } catch (error) {
    console.error('Failed to fetch tournaments:', error)
  } finally {
    loading.value = false
  }
}

function goToDetail(id) {
  router.push(`/tournaments/${id}`)
}

function watchLive(id) {
  router.push(`/tournaments/${id}`)
}

async function joinTournament(id) {
  try {
    await axios.post(`/api/v1/tournaments/${id}/join`)
    await fetchTournaments()
  } catch (error) {
    console.error('Failed to join tournament:', error)
  }
}

function getStatusBadge(status) {
  switch (status) {
    case 'LIVE': return 'hemingway-badge-success'
    case 'PLANNED': return 'hemingway-badge-info'
    case 'FINISHED': return 'hemingway-badge-warning'
    default: return 'hemingway-badge-info'
  }
}

function formatDate(dateString) {
  if (!dateString) return 'TBD'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.tournaments-view {
  max-width: 1200px;
  margin: 0 auto;
}

.tournaments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--hemmingway-space-xl);
}

.filters {
  display: flex;
  gap: var(--hemmingway-space-md);
  margin-bottom: var(--hemmingway-space-xl);
}

.filter-btn {
  padding: var(--hemmingway-space-sm) var(--hemmingway-space-lg);
  background: var(--hemmingway-bg-secondary);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--hemmingway-radius-full);
  color: var(--hemmingway-text-secondary);
  cursor: pointer;
  transition: all var(--hemmingway-transition-fast);
}

.filter-btn:hover {
  border-color: var(--hemmingway-turquoise);
}

.filter-btn.active {
  background: var(--hemmingway-neon);
  color: var(--hemmingway-dark);
  border-color: var(--hemmingway-neon);
}

.loading, .no-tournaments {
  text-align: center;
  padding: var(--hemmingway-space-2xl);
}

.tournaments-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--hemmingway-space-lg);
}

.tournament-card {
  cursor: pointer;
  transition: transform var(--hemmingway-transition-fast);
}

.tournament-card:hover {
  transform: translateY(-4px);
}

.tournament-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--hemmingway-space-lg);
}

.tournament-name {
  font-size: var(--hemmingway-font-size-lg);
  color: var(--hemmingway-text-primary);
  margin: 0;
}

.tournament-info {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-sm);
  margin-bottom: var(--hemmingway-space-lg);
}

.info-row {
  display: flex;
  justify-content: space-between;
}

.info-label {
  color: var(--hemmingway-text-muted);
  font-size: var(--hemmingway-font-size-sm);
}

.info-value {
  color: var(--hemmingway-text-primary);
  font-size: var(--hemmingway-font-size-sm);
}

.tournament-actions {
  display: flex;
  gap: var(--hemmingway-space-md);
}

.hemingway-btn-sm {
  padding: var(--hemmingway-space-sm) var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
}
</style>