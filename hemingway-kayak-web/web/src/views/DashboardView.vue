<template>
  <div class="dash">
    <header class="dash-head">
      <div class="dash-id">
        <span class="kicker"><i></i>Panel</span>
        <h1>Hola, <strong>{{ firstName }}</strong></h1>
        <p class="sub">{{ saludo }}</p>
      </div>
      <div class="dash-actions">
        <router-link to="/tracking" class="btn-primary">Empezar una salida</router-link>
        <router-link to="/history" class="btn-ghost">Ver historial</router-link>
      </div>
    </header>

    <div v-if="loading" class="state">Cargando tu panel…</div>

    <template v-else>
      <!-- Sin datos todavía: el panel no muestra ceros, invita a empezar -->
      <section v-if="isEmpty" class="empty">
        <img src="/logo.svg" alt="" class="empty-rose" />
        <h2>Todavía no hay salidas</h2>
        <p>
          Tu panel se llena solo: cada salida que grabes desde la app Android deja su
          recorrido, sus waypoints y sus capturas acá.
        </p>
        <div class="dash-actions center">
          <router-link to="/descargar" class="btn-primary">Instalar la app</router-link>
          <router-link to="/tournaments" class="btn-ghost">Ver torneos</router-link>
        </div>
      </section>

      <template v-else>
        <div class="kpis">
          <div class="kpi">
            <span class="kpi-label">Salidas</span>
            <b class="kpi-value">{{ stats.trips }}</b>
            <span class="kpi-sub">{{ stats.trips === 1 ? 'registrada' : 'registradas' }}</span>
          </div>
          <div class="kpi">
            <span class="kpi-label">Waypoints</span>
            <b class="kpi-value">{{ stats.spots }}</b>
            <span class="kpi-sub">puntos marcados</span>
          </div>
          <div class="kpi">
            <span class="kpi-label">Sesiones en vivo</span>
            <b class="kpi-value">{{ stats.sessions }}</b>
            <span class="kpi-sub">con tracking activo</span>
          </div>
          <div class="kpi">
            <span class="kpi-label">Torneos</span>
            <b class="kpi-value">{{ stats.tournaments }}</b>
            <span class="kpi-sub">participaciones</span>
          </div>
        </div>

        <div class="dash-grid">
          <section class="panel">
            <div class="panel-head">
              <h3>En el agua</h3>
              <router-link to="/tracking" class="panel-link">Abrir mapa</router-link>
            </div>
            <div v-if="!activeTrackers.length" class="panel-empty">Nadie está transmitiendo ahora.</div>
            <ul v-else class="water">
              <li v-for="t in activeTrackers" :key="t.userId">
                <span class="dot"></span>
                <div class="water-who">
                  <b>{{ trackerName(t) }}</b>
                  <span>hace {{ lastUpdate(t) }}</span>
                </div>
                <span class="water-pos">{{ formatLocation(t) }}</span>
              </li>
            </ul>
          </section>

          <section class="panel">
            <div class="panel-head">
              <h3>Próximos torneos</h3>
              <router-link to="/tournaments" class="panel-link">Todos</router-link>
            </div>
            <div v-if="!tournaments.length" class="panel-empty">No hay torneos programados.</div>
            <ul v-else class="tour">
              <li v-for="t in tournaments" :key="t.id">
                <div class="tour-when">
                  <b>{{ day(t.startTime) }}</b>
                  <span>{{ month(t.startTime) }}</span>
                </div>
                <div class="tour-info">
                  <b>{{ t.name }}</b>
                  <span>{{ t.locationName || 'Sin ubicación' }}</span>
                </div>
                <span class="badge" :class="statusClass(t.status)">{{ statusLabel(t.status) }}</span>
              </li>
            </ul>
          </section>
        </div>

        <section class="panel">
          <div class="panel-head">
            <h3>Últimas salidas</h3>
            <router-link to="/history" class="panel-link">Historial completo</router-link>
          </div>
          <div v-if="!trips.length" class="panel-empty">Sin salidas registradas.</div>
          <div v-else class="trips">
            <router-link
              v-for="trip in trips"
              :key="trip.id"
              :to="`/trips/${trip.id}`"
              class="trip"
            >
              <div class="trip-main">
                <b>{{ trip.title || 'Salida sin título' }}</b>
                <span class="trip-meta">
                  {{ formatDate(trip.startTime) }} · {{ trip.type || 'FISHING' }}
                </span>
              </div>
              <div class="trip-stats">
                <span v-if="trip.distance">{{ trip.distance.toFixed(1) }} km</span>
                <span v-else class="trip-live">EN VIVO</span>
              </div>
            </router-link>
          </div>
        </section>
      </template>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const loading = ref(true)
const stats = ref({ trips: 0, spots: 0, sessions: 0, tournaments: 0 })
const activeTrackers = ref([])
const tournaments = ref([])
const trips = ref([])

const firstName = computed(() => {
  const name = authStore.user?.displayName || authStore.user?.username || ''
  return name.split(/\s+/)[0] || 'kayaker'
})

const saludo = computed(() => {
  const n = stats.value.trips
  if (!n) return 'Grabá tu primera salida y el panel se encarga del resto.'
  if (n < 5) return `${n} salida${n === 1 ? '' : 's'} en el agua. Buen comienzo.`
  return `${n} salidas acumuladas. El mar te conoce bien.`
})

const isEmpty = computed(() =>
  !stats.value.trips && !stats.value.spots && !stats.value.tournaments && !stats.value.sessions
)

onMounted(load)

async function load() {
  loading.value = true
  // Salidas propias: /trips es el feed público y /trips/:id es del dueno,
  // por lo que linkear una salida ajena daba 403.
  const userId = authStore.user?.id
  const [statsRes, tripsRes, tourRes] = await Promise.allSettled([
    axios.get('/api/v1/users/me/stats'),
    userId
      ? axios.get(`/api/v1/history/user/${userId}/trips`)
      : Promise.resolve({ data: [] }),
    axios.get('/api/v1/tournaments')
  ])

  if (statsRes.status === 'fulfilled') stats.value = statsRes.value.data
  if (tripsRes.status === 'fulfilled') trips.value = (tripsRes.value.data || []).slice(0, 5)
  if (tourRes.status === 'fulfilled') {
    tournaments.value = (tourRes.value.data || [])
      .filter(t => t.status !== 'FINISHED' && t.status !== 'CANCELLED')
      .slice(0, 3)
  }

  try {
    const followers = await axios.get('/api/v1/tracking/followers')
    activeTrackers.value = followers.data || []
  } catch (e) {
    console.error('No se pudo cargar el tracking en vivo:', e)
  } finally {
    loading.value = false
  }
}

function trackerName(t) {
  return t.displayName || t.username || t.userId || 'Kayaker'
}

function lastUpdate(t) {
  const raw = t.timestamp || t.updatedAt
  if (!raw) return 'un rato'
  const time = new Date(raw)
  if (Number.isNaN(time.getTime())) return 'un rato'
  const mins = Math.max(0, Math.round((Date.now() - time.getTime()) / 60000))
  if (mins < 1) return 'ahora'
  if (mins < 60) return `${mins} min`
  return `${Math.round(mins / 60)} h`
}

function formatLocation(t) {
  if (t.latitude == null || t.longitude == null) return 'sin señal'
  return `${t.latitude.toFixed(3)}, ${t.longitude.toFixed(3)}`
}

function formatDate(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('es-CR', { day: '2-digit', month: 'short' })
}

function day(value) {
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? '—' : d.toLocaleDateString('es-CR', { day: '2-digit' })
}

function month(value) {
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? '' : d.toLocaleDateString('es-CR', { month: 'short' }).replace('.', '')
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
  if (status === 'UPCOMING' || status === 'REGISTRATION') return 'is-next'
  return 'is-done'
}
</script>

<style scoped>
.dash { padding: 48px 0 24px; }

.dash-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 32px;
  flex-wrap: wrap;
  padding-bottom: 32px;
  border-bottom: 1px solid var(--hemmingway-line);
}

.dash-head h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 44px;
  letter-spacing: -.02em;
  line-height: 1.05;
  margin: 14px 0 10px;
}

.dash-head h1 strong { font-weight: 600; color: var(--hemmingway-blue); }

.sub {
  color: var(--hemmingway-text-secondary);
  font-size: 16px;
  line-height: 1.55;
  max-width: 52ch;
}

.dash-actions { display: flex; gap: 12px; flex-wrap: wrap; }
.dash-actions.center { justify-content: center; }

.btn-primary {
  background: linear-gradient(135deg, var(--hemmingway-blue), var(--hemmingway-turquoise));
  color: var(--hemmingway-dark);
  font-size: 14px;
  font-weight: 700;
  padding: 12px 22px;
  border-radius: var(--hemmingway-radius-full);
  text-decoration: none;
  transition: all var(--hemmingway-transition-fast);
}

.btn-primary:hover { box-shadow: var(--hemmingway-shadow-glow); transform: translateY(-1px); }

.btn-ghost {
  border: 1px solid var(--hemmingway-line-2);
  color: var(--hemmingway-text-primary);
  font-size: 14px;
  font-weight: 600;
  padding: 12px 22px;
  border-radius: var(--hemmingway-radius-full);
  text-decoration: none;
  transition: all var(--hemmingway-transition-fast);
}

.btn-ghost:hover { border-color: var(--hemmingway-blue); }

.state { padding: 64px 0; color: var(--hemmingway-text-muted); font-size: 15px; }

.empty {
  text-align: center;
  padding: 64px 24px;
  max-width: 560px;
  margin: 0 auto;
}

.empty-rose { width: 84px; height: 84px; opacity: .9; margin-bottom: 24px; }

.empty h2 {
  font-family: var(--hemmingway-font-display);
  font-weight: 400;
  font-size: 30px;
  margin-bottom: 12px;
}

.empty p {
  color: var(--hemmingway-text-secondary);
  font-size: 15px;
  line-height: 1.65;
  margin-bottom: 28px;
}

.kpis {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 32px 0;
}

.kpi {
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 22px 22px 20px;
  background: linear-gradient(180deg, rgba(255,255,255,.02), transparent);
}

.kpi-label {
  display: block;
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
  margin-bottom: 14px;
}

.kpi-value {
  display: block;
  font-family: var(--hemmingway-font-display);
  font-weight: 500;
  font-size: 42px;
  line-height: 1;
  color: var(--hemmingway-text-primary);
  margin-bottom: 8px;
}

.kpi-sub { font-size: 12.5px; color: var(--hemmingway-text-secondary); }

.dash-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.panel {
  border: 1px solid var(--hemmingway-line);
  border-radius: var(--hemmingway-radius-lg);
  padding: 22px;
}

.panel + .panel { margin-top: 0; }

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
  letter-spacing: -.01em;
}

.panel-link {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  letter-spacing: .06em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
  text-decoration: none;
}

.panel-link:hover { color: var(--hemmingway-blue); }

.panel-empty {
  color: var(--hemmingway-text-muted);
  font-size: 14px;
  padding: 18px 0 6px;
}

.water { list-style: none; padding: 0; margin: 0; }

.water li {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 0;
  border-bottom: 1px solid var(--hemmingway-line);
}

.water li:last-child { border-bottom: none; }

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--hemmingway-neon);
  box-shadow: 0 0 10px var(--hemmingway-neon);
  flex: none;
}

.water-who { flex: 1; min-width: 0; }
.water-who b { display: block; font-size: 14.5px; font-weight: 600; }
.water-who span {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  color: var(--hemmingway-text-muted);
}

.water-pos {
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  color: var(--hemmingway-text-secondary);
}

.tour { list-style: none; padding: 0; margin: 0; }

.tour li {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid var(--hemmingway-line);
}

.tour li:last-child { border-bottom: none; }

.tour-when {
  flex: none;
  width: 52px;
  text-align: center;
  border: 1px solid var(--hemmingway-line-2);
  border-radius: var(--hemmingway-radius-md);
  padding: 6px 0;
}

.tour-when b {
  display: block;
  font-family: var(--hemmingway-font-display);
  font-size: 20px;
  font-weight: 500;
  line-height: 1;
}

.tour-when span {
  font-family: var(--hemmingway-font-mono);
  font-size: 10px;
  letter-spacing: .08em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
}

.tour-info { flex: 1; min-width: 0; }
.tour-info b { display: block; font-size: 14.5px; font-weight: 600; }
.tour-info span {
  font-size: 12.5px;
  color: var(--hemmingway-text-secondary);
}

.badge {
  flex: none;
  font-family: var(--hemmingway-font-mono);
  font-size: 10px;
  letter-spacing: .1em;
  padding: 5px 9px;
  border-radius: var(--hemmingway-radius-full);
  border: 1px solid var(--hemmingway-line-2);
  color: var(--hemmingway-text-secondary);
}

.badge.is-live {
  border-color: var(--hemmingway-neon);
  color: var(--hemmingway-neon);
}

.badge.is-next { color: var(--hemmingway-blue); border-color: var(--hemmingway-blue); }
.badge.is-done { color: var(--hemmingway-text-muted); }

.trips { display: flex; flex-direction: column; }

.trip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 14px 0;
  border-bottom: 1px solid var(--hemmingway-line);
  text-decoration: none;
  color: inherit;
  transition: padding var(--hemmingway-transition-fast);
}

.trip:last-child { border-bottom: none; }

.trip:hover { padding-left: 6px; }

.trip-main b {
  display: block;
  font-size: 15px;
  font-weight: 600;
  color: var(--hemmingway-text-primary);
}

.trip-meta {
  font-family: var(--hemmingway-font-mono);
  font-size: 11.5px;
  color: var(--hemmingway-text-muted);
  letter-spacing: .04em;
}

.trip-stats {
  flex: none;
  font-family: var(--hemmingway-font-mono);
  font-size: 13px;
  color: var(--hemmingway-gold);
}

.trip-live { color: var(--hemmingway-neon); }

@media (max-width: 900px) {
  .dash { padding: 32px 0 24px; }
  .dash-head h1 { font-size: 32px; }
  .kpis { grid-template-columns: repeat(2, 1fr); }
  .dash-grid { grid-template-columns: 1fr; }
}

@media (max-width: 520px) {
  .kpis { grid-template-columns: 1fr; }
}
</style>
