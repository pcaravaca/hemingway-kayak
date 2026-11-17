<template>
  <div class="app-shell">
    <header class="nav">
      <router-link to="/" class="brand">
        <img src="/logo.svg" alt="Hemingway Kayak" class="brand-mark" />
        <span class="brand-word">Hemingway<em>Kayak</em></span>
      </router-link>

      <nav class="nav-links">
        <router-link v-if="isAuthenticated" to="/dashboard">Panel</router-link>
        <router-link to="/tracking">Tracking</router-link>
        <router-link to="/tournaments">Torneos</router-link>
        <router-link to="/history">Historial</router-link>
        <router-link to="/descargar">Descargar</router-link>
      </nav>

      <div class="nav-actions">
        <template v-if="isAuthenticated">
          <router-link to="/profile" class="btn-ghost">Mi perfil</router-link>
          <button class="btn-ghost" type="button" @click="handleLogout">Salir</button>
        </template>
        <template v-else>
          <router-link to="/login" class="btn-ghost">Entrar</router-link>
          <router-link to="/register" class="btn-solid">Crear cuenta</router-link>
        </template>
      </div>
    </header>

    <main class="main-content">
      <router-view />
    </main>

    <footer class="foot">
      <span>HEMINGWAY KAYAK · COSTA RICA</span>
      <span>9.9382° N — 84.8371° W</span>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const isAuthenticated = computed(() => authStore.isAuthenticated)

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style>
/* La retícula de carta cubre todo el shell */
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

.app-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background-image:
    linear-gradient(var(--hemmingway-line) 1px, transparent 1px),
    linear-gradient(90deg, var(--hemmingway-line) 1px, transparent 1px);
  background-size: 64px 64px;
  opacity: .35;
  -webkit-mask-image: radial-gradient(120% 80% at 50% 0%, #000 0%, transparent 70%);
          mask-image: radial-gradient(120% 80% at 50% 0%, #000 0%, transparent 70%);
}

.nav, .main-content, .foot { position: relative; z-index: 1; }
</style>

<style scoped>
.nav {
  display: flex;
  align-items: center;
  gap: 32px;
  padding: 20px 40px;
  border-bottom: 1px solid var(--hemmingway-line);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: var(--hemmingway-text-primary);
}

.brand-mark { width: 34px; height: 34px; display: block; }

.brand-word {
  font-family: var(--hemmingway-font-display);
  font-size: 20px;
  font-weight: 500;
  letter-spacing: -.01em;
}

.brand-word em { font-style: normal; color: var(--hemmingway-blue); }

.nav-links { display: flex; gap: 26px; margin-left: auto; }

.nav-links a {
  color: var(--hemmingway-text-secondary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color var(--hemmingway-transition-fast);
}

.nav-links a:hover,
.nav-links a.router-link-active { color: var(--hemmingway-text-primary); }

.nav-actions { display: flex; gap: 10px; align-items: center; }

.btn-ghost {
  border: 1px solid var(--hemmingway-line-2);
  color: var(--hemmingway-text-primary);
  background: transparent;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  padding: 9px 16px;
  border-radius: var(--hemmingway-radius-full);
  text-decoration: none;
  cursor: pointer;
  transition: all var(--hemmingway-transition-fast);
}

.btn-ghost:hover { border-color: var(--hemmingway-blue); }

.btn-solid {
  background: var(--hemmingway-blue);
  color: var(--hemmingway-dark);
  font-size: 13px;
  font-weight: 700;
  padding: 10px 18px;
  border-radius: var(--hemmingway-radius-full);
  text-decoration: none;
  transition: all var(--hemmingway-transition-fast);
}

.btn-solid:hover { box-shadow: var(--hemmingway-shadow-glow); }

.main-content {
  flex: 1;
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding: 0 40px;
}

@media (max-width: 900px) {
  .nav { flex-wrap: wrap; gap: 16px; padding: 16px 20px; }
  .nav-links { margin-left: 0; order: 3; width: 100%; }
  .nav-actions { margin-left: auto; }
  .main-content { padding: 0 20px; }
}
</style>
