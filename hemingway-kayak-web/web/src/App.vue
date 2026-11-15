<template>
  <div id="app" class="hemingway-body">
    <nav class="navbar">
      <div class="nav-brand">
        <router-link to="/" class="brand-link">
          <img src="/logo.svg" alt="Hemingway Kayak" class="logo" />
          <span class="brand-text">Hemingway Kayak</span>
        </router-link>
      </div>
      <div class="nav-links">
        <router-link to="/">Home</router-link>
        <router-link to="/tracking">Tracking</router-link>
        <router-link to="/tournaments">Tournaments</router-link>
        <router-link to="/history">History</router-link>
        <router-link v-if="isAuthenticated" to="/profile">Profile</router-link>
        <router-link v-else to="/login">Login</router-link>
      </div>
    </nav>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>

<style>
@import './src/assets/hemingway.css';

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 32px;
  background: var(--hemmingway-bg-secondary);
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo {
  width: 40px;
  height: 40px;
}

.brand-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--hemmingway-neon);
}

.nav-links {
  display: flex;
  gap: 24px;
}

.nav-links a {
  color: var(--hemmingway-text-secondary);
  text-decoration: none;
  font-weight: 500;
  transition: color 150ms ease;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: var(--hemmingway-neon);
}

.brand-link {
  text-decoration: none;
}

.main-content {
  min-height: calc(100vh - 72px);
  padding: 32px;
}
</style>
