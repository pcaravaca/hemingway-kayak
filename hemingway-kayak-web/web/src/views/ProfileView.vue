<template>
  <div class="profile-view">
    <div class="profile-header">
      <h1 class="hemingway-title">Perfil</h1>
    </div>

    <div v-if="loading" class="loading">Cargando perfil...</div>

    <div v-else-if="!isAuthenticated" class="not-authenticated">
      <div class="empty-state">
        <div class="empty-icon">🔐</div>
        <h3>Tenés que entrar</h3>
        <p class="text-muted">Iniciá sesión para ver tu perfil</p>
        <router-link to="/login" class="hemingway-btn hemingway-btn-primary">
          Sign In
        </router-link>
      </div>
    </div>

    <div v-else class="profile-content">
      <div class="profile-card hemingway-card">
        <div class="profile-avatar">
          <img
            v-if="user?.avatarUrl"
            :src="user.avatarUrl"
            class="avatar-img"
            alt="Foto de perfil"
          />
          <div v-else class="avatar-placeholder">
            {{ userInitials }}
          </div>
        </div>

        <div class="profile-info">
          <h2 class="profile-name">{{ user?.displayName || user?.username }}</h2>
          <p class="profile-email text-muted">{{ user?.email }}</p>
        </div>

        <div class="profile-stats">
          <div class="stat-item">
            <span class="stat-value">{{ stats.trips }}</span>
            <span class="stat-label">Trips</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ stats.tournaments }}</span>
            <span class="stat-label">Tournaments</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ stats.catches }}</span>
            <span class="stat-label">Catches</span>
          </div>
        </div>

        <div class="profile-settings">
          <h3 class="settings-title">Configuración</h3>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">Tracking en vivo</span>
              <span class="setting-desc text-muted">Compartí tu ubicación con otros</span>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="trackingEnabled" @change="updateTrackingSetting" />
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>

        <div class="profile-actions">
          <button class="hemingway-btn hemingway-btn-secondary" @click="editProfile">
            Edit Profile
          </button>
          <button class="hemingway-btn hemingway-btn-ghost" @click="handleLogout">
            Sign Out
          </button>
        </div>
      </div>

      <div class="kayaks-section">
        <h2 class="section-title">Mis kayaks</h2>

        <div v-if="kayaks.length === 0" class="no-kayaks">
          <p class="text-muted">Todavía no hay kayaks</p>
        </div>

        <div v-else class="kayaks-grid">
          <div v-for="kayak in kayaks" :key="kayak.id" class="kayak-card hemingway-card">
            <h4 class="kayak-name">{{ kayak.name }}</h4>
            <div class="kayak-details">
              <span>{{ kayak.type }}</span>
              <span v-if="kayak.length">{{ kayak.length }}cm</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const kayaks = ref([])
const trackingEnabled = ref(false)
const stats = ref({ trips: 0, tournaments: 0, catches: 0 })

const isAuthenticated = computed(() => authStore.isAuthenticated)
const user = computed(() => authStore.user)
const userInitials = computed(() => {
  const name = user.value?.displayName || user.value?.username || ''
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

onMounted(async () => {
  if (isAuthenticated.value) {
    await Promise.all([fetchKayaks(), fetchStats()])
    trackingEnabled.value = user.value?.trackingEnabled || false
  }
  loading.value = false
})

async function fetchKayaks() {
  try {
    const response = await axios.get(`/api/v1/kayaks?userId=${user.value?.id}`)
    kayaks.value = response.data
  } catch (error) {
    console.error('Failed to fetch kayaks:', error)
  }
}

async function fetchStats() {
  try {
    const response = await axios.get('/api/v1/users/me/stats')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

async function updateTrackingSetting() {
  try {
    await axios.put(`/api/v1/users/${user.value?.id}`, {
      trackingEnabled: trackingEnabled.value
    })
  } catch (error) {
    console.error('Failed to update tracking setting:', error)
    trackingEnabled.value = !trackingEnabled.value
  }
}

function editProfile() {
  router.push('/profile/edit')
}

async function handleLogout() {
  authStore.logout()
  router.push('/')
}
</script>

<style scoped>
.profile-view {
  max-width: 800px;
  margin: 0 auto;
}

.profile-header {
  margin-bottom: var(--hemmingway-space-xl);
}

.loading, .not-authenticated {
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
  margin-bottom: var(--hemmingway-space-sm);
}

.empty-state .hemingway-btn {
  margin-top: var(--hemmingway-space-lg);
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-xl);
}

.profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.profile-avatar {
  margin-bottom: var(--hemmingway-space-lg);
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--hemmingway-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--hemmingway-font-size-2xl);
  font-weight: 700;
  color: var(--hemmingway-dark);
}

.avatar-img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.profile-name {
  font-size: var(--hemmingway-font-size-2xl);
  color: var(--hemmingway-text-primary);
  margin: 0;
}

.profile-email {
  margin-top: var(--hemmingway-space-xs);
}

.profile-stats {
  display: flex;
  gap: var(--hemmingway-space-2xl);
  margin: var(--hemmingway-space-xl) 0;
  padding: var(--hemmingway-space-lg) 0;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: var(--hemmingway-font-size-2xl);
  font-weight: 700;
  color: var(--hemmingway-neon);
}

.stat-label {
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-text-muted);
}

.profile-settings {
  width: 100%;
  text-align: left;
  margin-bottom: var(--hemmingway-space-xl);
}

.settings-title {
  font-size: var(--hemmingway-font-size-lg);
  color: var(--hemmingway-text-primary);
  margin-bottom: var(--hemmingway-space-lg);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--hemmingway-space-md) 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.setting-label {
  font-weight: 500;
  color: var(--hemmingway-text-primary);
}

.setting-desc {
  font-size: var(--hemmingway-font-size-sm);
}

.toggle {
  position: relative;
  width: 48px;
  height: 24px;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background: var(--hemmingway-bg-tertiary);
  border-radius: var(--hemmingway-radius-full);
  transition: var(--hemmingway-transition-fast);
}

.toggle-slider::before {
  content: '';
  position: absolute;
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background: var(--hemmingway-text-primary);
  border-radius: 50%;
  transition: var(--hemmingway-transition-fast);
}

.toggle input:checked + .toggle-slider {
  background: var(--hemmingway-neon);
}

.toggle input:checked + .toggle-slider::before {
  transform: translateX(24px);
  background: var(--hemmingway-dark);
}

.profile-actions {
  display: flex;
  gap: var(--hemmingway-space-md);
}

.section-title {
  font-size: var(--hemmingway-font-size-xl);
  color: var(--hemmingway-neon);
  margin-bottom: var(--hemmingway-space-lg);
}

.kayaks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--hemmingway-space-md);
}

.kayak-card {
  text-align: left;
}

.kayak-name {
  font-size: var(--hemmingway-font-size-base);
  color: var(--hemmingway-text-primary);
  margin: 0 0 var(--hemmingway-space-sm) 0;
}

.kayak-details {
  display: flex;
  gap: var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
  color: var(--hemmingway-text-muted);
}
</style>