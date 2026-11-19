<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const form = ref({
  title: '',
  description: '',
  type: 'FISHING',
  startTime: '',
  startLat: '',
  startLon: ''
})
const loading = ref(false)
const error = ref(null)
const locationError = ref(null)

const TYPES = ['FISHING', 'CAMPING', 'TOUR', 'RACING']

async function useMyLocation() {
  locationError.value = null
  if (!navigator.geolocation) {
    locationError.value = 'Este navegador no tiene geolocalización'
    return
  }
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      form.value.startLat = pos.coords.latitude.toFixed(6)
      form.value.startLon = pos.coords.longitude.toFixed(6)
    },
    () => {
      locationError.value = 'No se pudo obtener tu ubicación'
    }
  )
}

async function save() {
  error.value = null
  if (!form.value.title.trim()) {
    error.value = 'El título es obligatorio'
    return
  }
  loading.value = true
  try {
    await axios.post('/api/v1/trips/', {
      title: form.value.title.trim(),
      description: form.value.description || null,
      type: form.value.type,
      startLat: form.value.startLat ? Number(form.value.startLat) : null,
      startLon: form.value.startLon ? Number(form.value.startLon) : null
    })
    router.push('/history')
  } catch (e) {
    error.value = e.response?.data?.error || 'No se pudo crear la salida'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="trip-create">
    <div class="page-head">
      <div>
        <span class="kicker">Bitácora</span>
        <h1 class="hemmingway-title">Nueva salida</h1>
      </div>
      <button class="hemmingway-btn hemmingway-btn-ghost" @click="router.push('/history')">
        Volver
      </button>
    </div>

    <form class="hemmingway-card form" @submit.prevent="save">
      <label class="field">
        <span>Título</span>
        <input v-model="form.title" class="hemmingway-input" required />
      </label>

      <label class="field">
        <span>Descripción</span>
        <textarea v-model="form.description" class="hemmingway-input" rows="3" />
      </label>

      <label class="field">
        <span>Tipo</span>
        <select v-model="form.type" class="hemmingway-input">
          <option v-for="t in TYPES" :key="t" :value="t">{{ t }}</option>
        </select>
      </label>

      <label class="field">
        <span>Fecha de inicio</span>
        <input v-model="form.startTime" type="datetime-local" class="hemmingway-input" />
      </label>

      <div class="field-row">
        <label class="field">
          <span>Latitud</span>
          <input v-model="form.startLat" type="number" step="any" class="hemmingway-input" placeholder="9.9382" />
        </label>
        <label class="field">
          <span>Longitud</span>
          <input v-model="form.startLon" type="number" step="any" class="hemmingway-input" placeholder="-84.8371" />
        </label>
      </div>

      <button type="button" class="hemmingway-btn hemmingway-btn-ghost" @click="useMyLocation">
        Usar mi ubicación
      </button>
      <p v-if="locationError" class="form-error">{{ locationError }}</p>
      <p v-if="error" class="form-error">{{ error }}</p>

      <div class="form-actions">
        <button type="submit" class="hemmingway-btn hemmingway-btn-primary" :disabled="loading">
          {{ loading ? 'Guardando...' : 'Crear salida' }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 560px;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.field > span {
  font-size: 0.8rem;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
}
.form-error {
  margin: 0;
  color: #ff8a80;
  font-size: 0.85rem;
}
</style>
