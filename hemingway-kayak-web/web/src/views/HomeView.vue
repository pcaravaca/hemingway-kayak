<template>
  <div class="home">
    <header class="hero">
      <h1 class="hemingway-title">Hemingway Kayak</h1>
      <p class="hemingway-subtitle">Track your fishing adventures in real-time</p>
    </header>

    <section class="features">
      <div class="feature-card hemingway-card">
        <div class="feature-icon">🗺️</div>
        <h3>Live Tracking</h3>
        <p>Share your location with friends and family while on the water</p>
      </div>

      <div class="feature-card hemingway-card">
        <div class="feature-icon">🏆</div>
        <h3>Tournaments</h3>
        <p>Join live fishing tournaments with real-time rankings</p>
      </div>

      <div class="feature-card hemingway-card">
        <div class="feature-icon">📊</div>
        <h3>History</h3>
        <p>Review your past trips, catches, and achievements</p>
      </div>
    </section>

    <section class="live-trackers" v-if="activeTrackers.length > 0">
      <h2 class="section-title">Live Now</h2>
      <div class="trackers-list">
        <div v-for="tracker in activeTrackers" :key="tracker.userId" class="tracker-item hemingway-card">
          <div class="tracker-info">
            <span class="tracker-name">{{ tracker.displayName }}</span>
            <span class="tracker-status hemingway-badge hemingway-badge-success">LIVE</span>
          </div>
          <div class="tracker-location">
            {{ formatLocation(tracker) }}
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const activeTrackers = ref([])

onMounted(async () => {
  try {
    const response = await axios.get('/api/v1/tracking/followers')
    activeTrackers.value = response.data.filter(t => t.isActive)
  } catch (error) {
    console.error('Failed to load trackers:', error)
  }
})

function formatLocation(tracker) {
  return `${tracker.latitude.toFixed(4)}, ${tracker.longitude.toFixed(4)}`
}
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
}

.hero {
  text-align: center;
  padding: var(--hemmingway-space-2xl) 0;
}

.section-title {
  font-size: var(--hemmingway-font-size-2xl);
  color: var(--hemmingway-neon);
  margin-bottom: var(--hemmingway-space-lg);
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--hemmingway-space-lg);
  margin-bottom: var(--hemmingway-space-2xl);
}

.feature-card {
  text-align: center;
  transition: transform var(--hemmingway-transition-normal);
}

.feature-card:hover {
  transform: translateY(-4px);
}

.feature-icon {
  font-size: 48px;
  margin-bottom: var(--hemmingway-space-md);
}

.live-trackers {
  margin-top: var(--hemmingway-space-2xl);
}

.tracker-item {
  margin-bottom: var(--hemmingway-space-md);
}

.tracker-info {
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
</style>
