<template>
  <div class="create">
    <router-link to="/tournaments" class="back">← Volver a torneos</router-link>

    <header class="create-head">
      <span class="kicker"><i></i>Torneos</span>
      <h1>Crear un torneo</h1>
      <p class="sub">Quedás como organizador; los demás kayakers se suman desde la lista.</p>
    </header>

    <form class="card form" @submit.prevent="handleCreate">
      <label class="field">
        <span>Nombre</span>
        <input v-model="form.name" class="hemmingway-input" placeholder="Ej.: Copa Golfo de Nicoya" required />
      </label>

      <label class="field">
        <span>Descripción</span>
        <textarea
          v-model="form.description"
          class="hemmingway-input"
          rows="3"
          placeholder="Reglas, premios, punto de encuentro"
        ></textarea>
      </label>

      <div class="row">
        <label class="field">
          <span>Tipo</span>
          <select v-model="form.type" class="hemmingway-input">
            <option value="FISHING">Pesca</option>
            <option value="RACING">Carrera</option>
            <option value="TOURING">Travesía</option>
          </select>
        </label>
        <label class="field">
          <span>Cupo máximo</span>
          <input v-model.number="form.maxParticipants" type="number" min="2" class="hemmingway-input" />
        </label>
      </div>

      <div class="row">
        <label class="field">
          <span>Empieza</span>
          <input v-model="form.startTime" type="datetime-local" class="hemmingway-input" required />
        </label>
        <label class="field">
          <span>Termina</span>
          <input v-model="form.endTime" type="datetime-local" class="hemmingway-input" required />
        </label>
      </div>

      <label class="field">
        <span>Ubicación</span>
        <input v-model="form.locationName" class="hemmingway-input" placeholder="Ej.: Playa Herradura" />
      </label>

      <div class="row">
        <label class="field">
          <span>Latitud (opcional)</span>
          <input v-model.number="form.startLat" type="number" step="any" class="hemmingway-input" placeholder="9.9382" />
        </label>
        <label class="field">
          <span>Longitud (opcional)</span>
          <input v-model.number="form.startLon" type="number" step="any" class="hemmingway-input" placeholder="-84.8371" />
        </label>
      </div>

      <p v-if="message" class="message" :class="{ ok: created }">{{ message }}</p>

      <div class="actions">
        <button type="submit" class="hemmingway-btn hemmingway-btn-primary" :disabled="saving">
          {{ saving ? 'Creando…' : 'Crear torneo' }}
        </button>
        <router-link to="/tournaments" class="hemmingway-btn hemmingway-btn-secondary">Cancelar</router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const form = ref({
  name: '',
  description: '',
  type: 'FISHING',
  maxParticipants: 20,
  startTime: '',
  endTime: '',
  locationName: '',
  startLat: null,
  startLon: null
})

const saving = ref(false)
const created = ref(false)
const message = ref('')

async function handleCreate() {
  saving.value = true
  message.value = ''
  try {
    const payload = {
      name: form.value.name,
      description: form.value.description,
      type: form.value.type,
      maxParticipants: Number(form.value.maxParticipants) || null,
      startTime: new Date(form.value.startTime).toISOString(),
      endTime: new Date(form.value.endTime).toISOString(),
      locationName: form.value.locationName || null,
      startLat: form.value.startLat,
      startLon: form.value.startLon
    }
    await axios.post('/api/v1/tournaments', payload)
    created.value = true
    message.value = 'Torneo creado.'
    setTimeout(() => router.push('/tournaments'), 700)
  } catch (e) {
    created.value = false
    message.value = e.response?.data?.error || 'No se pudo crear el torneo.'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.create { padding: 40px 0 24px; max-width: 680px; }

.back {
  display: inline-block;
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  letter-spacing: .06em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
  text-decoration: none;
  margin-bottom: 22px;
}

.back:hover { color: var(--hemmingway-blue); }

.create-head { padding-bottom: 24px; }

.create-head h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 36px;
  letter-spacing: -.02em;
  margin: 12px 0 8px;
}

.sub { color: var(--hemmingway-text-secondary); font-size: 15px; }

.form { display: flex; flex-direction: column; gap: 18px; }

.row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

.field { display: flex; flex-direction: column; gap: 8px; }

.field span {
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
}

select.hemmingway-input,
input[type="datetime-local"].hemmingway-input { color-scheme: dark; }

.message { font-size: 14px; color: var(--hemmingway-accent-error); margin: 0; }
.message.ok { color: var(--hemmingway-accent-success); }

.actions { display: flex; gap: 12px; flex-wrap: wrap; }

@media (max-width: 620px) {
  .create { padding: 28px 0 24px; }
  .row { grid-template-columns: 1fr; }
}
</style>
