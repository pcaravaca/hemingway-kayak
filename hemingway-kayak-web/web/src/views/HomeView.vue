<template>
  <div class="home">
    <section class="hero">
      <div class="hero-copy">
        <span class="kicker"><i></i>{{ liveLabel }}</span>
        <h1>La pesca en kayak, <strong>medida</strong> hasta la última milla.</h1>
        <p class="lede">
          Compartí tu posición con los tuyos, marcá waypoints en los mejores puntos y
          reviví cada salida. Hemingway Kayak deja el registro donde el celular no llega.
        </p>
        <div class="cta">
          <router-link v-if="isAuthenticated" to="/dashboard" class="btn-solid">Ir a mi panel</router-link>
          <router-link v-else to="/register" class="btn-solid">Empezar una salida</router-link>
          <router-link to="/tournaments" class="btn-outline">Ver torneos</router-link>
        </div>
        <div class="coords">
          <div>EN EL AGUA<b>{{ activeTrackers.length }} kayaks</b></div>
          <div>WAYPOINTS<b>14</b></div>
          <div>SALIDAS<b>37</b></div>
        </div>
      </div>

      <div class="hero-art">
        <div class="ring r1"></div>
        <div class="ring r2"></div>
        <div class="ring r3"></div>
        <img class="rose" src="/logo.svg" alt="" />
        <div class="chip a">RUMBO<b>214° NE</b></div>
        <div class="chip b">VELOCIDAD<b>4.2 kn</b></div>
        <div class="chip c">PROFUNDIDAD<b>18 m</b></div>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <span>01 — Qué hace</span>
        <h2>Tres cosas, bien hechas</h2>
      </div>
      <div class="feature-grid">
        <div class="feature">
          <img class="icon" src="/icons/map-pin.svg" alt="" />
          <div class="num">01</div>
          <h3>Tracking en vivo</h3>
          <p>Tu posición se publica cada pocos segundos. Familia y amigos la siguen desde
             cualquier navegador, sin instalar nada.</p>
        </div>
        <div class="feature">
          <img class="icon" src="/icons/fish.svg" alt="" />
          <div class="num">02</div>
          <h3>Torneos</h3>
          <p>Sumate a competencias con ranking en tiempo real, capturas verificadas y
             premios definidos por el organizador.</p>
        </div>
        <div class="feature">
          <img class="icon" src="/icons/waves.svg" alt="" />
          <div class="num">03</div>
          <h3>Historial</h3>
          <p>Cada salida queda guardada con su recorrido, waypoints, capturas y condiciones.
             Revivila cuando quieras.</p>
        </div>
      </div>
    </section>

    <section class="section" v-if="activeTrackers.length">
      <div class="section-head">
        <span>02 — Ahora mismo</span>
        <h2>En el agua</h2>
      </div>
      <div class="live">
        <div v-for="tracker in activeTrackers" :key="tracker.userId" class="live-row">
          <div class="who">
            <div class="avatar">{{ initials(tracker) }}</div>
            <div>
              <b>{{ trackerName(tracker) }}</b>
              <span>hace {{ lastUpdate(tracker) }}</span>
            </div>
          </div>
          <span class="hemmingway-badge hemmingway-badge-success">EN VIVO</span>
          <span class="pos">{{ formatLocation(tracker) }}</span>
        </div>
      </div>
    </section>

    <!-- Sin sesión no se puede ver quién está en el agua (el endpoint es
         protegido): se ofrece el login en lugar de una sección vacía. -->
    <section class="section" v-else-if="!isAuthenticated">
      <div class="section-head">
        <span>02 — Ahora mismo</span>
        <h2>En el agua</h2>
      </div>
      <div class="live">
        <div class="live-row">
          <div class="who">
            <div>
              <b>Seguí a los kayakers en vivo</b>
              <span>La posición se comparte con tu cuenta</span>
            </div>
          </div>
          <router-link to="/login" class="hemmingway-btn hemmingway-btn-secondary">Entrar</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
const activeTrackers = ref([])

const liveLabel = computed(() =>
  activeTrackers.value.length
    ? `Tracking en vivo · ${activeTrackers.value.length} en el agua`
    : 'Tracking en vivo · Golfo de Nicoya'
)

async function loadTrackers() {
  // /tracking/followers es un endpoint protegido. Sin sesión no se pide: antes
  // se pedía igual y la consola se llenaba de 401 con la sección vacía.
  if (!authStore.isAuthenticated) {
    activeTrackers.value = []
    return
  }
  try {
    // El backend ya filtra is_active = TRUE.
    const response = await axios.get('/api/v1/tracking/followers')
    activeTrackers.value = response.data || []
  } catch (error) {
    console.error('No se pudo cargar el tracking en vivo:', error)
    activeTrackers.value = []
  }
}

onMounted(loadTrackers)
watch(isAuthenticated, loadTrackers)

function trackerName(tracker) {
  return tracker.displayName || tracker.username || tracker.userId || 'Kayaker'
}

function initials(tracker) {
  const name = trackerName(tracker)
  return name.split(/\s+/).slice(0, 2).map(w => w[0]).join('').toUpperCase().slice(0, 2)
}

function lastUpdate(tracker) {
  const raw = tracker.updatedAt || tracker.lastUpdate
  if (!raw) return 'un rato'
  const t = new Date(raw)
  if (Number.isNaN(t.getTime())) return 'un rato'
  const mins = Math.max(0, Math.round((Date.now() - t.getTime()) / 60000))
  if (mins < 1) return 'ahora'
  if (mins < 60) return `${mins} min`
  return `${Math.round(mins / 60)} h ${mins % 60} min`
}

function formatLocation(tracker) {
  if (tracker.latitude == null || tracker.longitude == null) return 'sin señal'
  return `${tracker.latitude.toFixed(4)}, ${tracker.longitude.toFixed(4)}`
}
</script>

<style scoped>
.home { padding-bottom: 20px; }

.hero {
  display: grid;
  grid-template-columns: 1.15fr .85fr;
  gap: 64px;
  align-items: center;
  padding: 72px 0 64px;
}

.hero h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 64px;
  line-height: 1.03;
  letter-spacing: -.02em;
  margin: 26px 0 22px;
  max-width: 13ch;
}

.hero h1 strong { font-weight: 600; color: var(--hemmingway-blue); }

.lede {
  font-size: 17px;
  line-height: 1.65;
  color: var(--hemmingway-text-secondary);
  max-width: 46ch;
}

.cta { display: flex; gap: 14px; margin-top: 34px; }

.btn-solid {
  background: var(--hemmingway-blue);
  color: var(--hemmingway-dark);
  font-size: 15px;
  font-weight: 700;
  padding: 14px 26px;
  border-radius: var(--hemmingway-radius-lg);
  text-decoration: none;
  transition: all var(--hemmingway-transition-normal);
}

.btn-solid:hover { box-shadow: var(--hemmingway-shadow-glow); transform: translateY(-2px); }

.btn-outline {
  border: 1px solid var(--hemmingway-line-2);
  color: var(--hemmingway-text-primary);
  font-size: 15px;
  font-weight: 600;
  padding: 14px 26px;
  border-radius: var(--hemmingway-radius-lg);
  text-decoration: none;
  transition: all var(--hemmingway-transition-normal);
}

.btn-outline:hover { border-color: var(--hemmingway-blue); }

.coords {
  display: flex;
  gap: 32px;
  margin-top: 44px;
  padding-top: 26px;
  border-top: 1px solid var(--hemmingway-line);
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  color: var(--hemmingway-text-muted);
}

.coords b {
  display: block;
  color: var(--hemmingway-gold);
  font-weight: 500;
  margin-top: 6px;
  font-size: 13px;
}

.hero-art { position: relative; height: 460px; }

.hero-art .ring {
  position: absolute;
  inset: 0;
  margin: auto;
  border: 1px solid var(--hemmingway-line);
  border-radius: 50%;
}

.hero-art .r1 { width: 440px; height: 440px; }
.hero-art .r2 { width: 320px; height: 320px; border-color: var(--hemmingway-line-2); }
.hero-art .r3 { width: 200px; height: 200px; border-style: dashed; }

.hero-art .rose {
  position: absolute;
  inset: 0;
  margin: auto;
  width: 150px;
  height: 150px;
}

.chip.a { top: 36px; left: -8px; }
.chip.b { bottom: 52px; right: -4px; }
.chip.c { top: 50%; left: -24px; }

.section { padding: 56px 0; border-top: 1px solid var(--hemmingway-line); }

@media (max-width: 980px) {
  .hero { grid-template-columns: 1fr; gap: 40px; padding: 40px 0; }
  .hero h1 { font-size: 42px; }
  .hero-art { height: 360px; }
  .hero-art .r1 { width: 320px; height: 320px; }
  .hero-art .r2 { width: 230px; height: 230px; }
  .hero-art .r3 { width: 150px; height: 150px; }
  .chip.c, .chip.a { left: 0; }
  .coords { flex-wrap: wrap; gap: 20px; }
}
</style>
