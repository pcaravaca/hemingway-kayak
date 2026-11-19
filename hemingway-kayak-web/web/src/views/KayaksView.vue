<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const kayaks = ref([])
const loading = ref(true)
const showForm = ref(false)
const editingId = ref(null)
const saving = ref(false)
const error = ref(null)

const emptyForm = {
  name: '',
  type: 'SIT_ON_TOP',
  length: '',
  width: '',
  weight: '',
  maxCapacity: '',
  material: 'FIBERGLASS',
  propulsion: 'PADDLE',
  imageUrl: ''
}

const form = ref({ ...emptyForm })
const confirmingId = ref(null)

const TYPES = ['SIT_ON_TOP', 'SIT_INSIDE', 'INFLATABLE', 'FISHING', 'RACING', 'TOURING']
const MATERIALS = ['FIBERGLASS', 'PLASTIC', 'CARBON', 'WOOD', 'INFLATABLE']
const PROPULSIONS = ['PADDLE', 'PEDAL', 'SAIL', 'MOTOR', 'MIXED']

async function load() {
  loading.value = true
  try {
    const { data } = await axios.get('/api/v1/kayaks/')
    kayaks.value = data
  } catch (e) {
    error.value = e.response?.data?.error || 'No se pudieron cargar los kayaks'
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { ...emptyForm }
  error.value = null
  showForm.value = true
}

function openEdit(kayak) {
  editingId.value = kayak.id
  form.value = {
    name: kayak.name || '',
    type: kayak.type || 'SIT_ON_TOP',
    length: kayak.length ?? '',
    width: kayak.width ?? '',
    weight: kayak.weight ?? '',
    maxCapacity: kayak.maxCapacity ?? '',
    material: kayak.material || 'FIBERGLASS',
    propulsion: kayak.propulsion || 'PADDLE',
    imageUrl: kayak.imageUrl || ''
  }
  error.value = null
  showForm.value = true
}

function closeForm() {
  showForm.value = false
}

async function save() {
  if (!form.value.name.trim()) {
    error.value = 'El nombre es obligatorio'
    return
  }
  saving.value = true
  error.value = null
  try {
    const payload = {
      name: form.value.name.trim(),
      type: form.value.type,
      length: Number(form.value.length) || 0,
      width: Number(form.value.width) || 0,
      weight: Number(form.value.weight) || 0,
      maxCapacity: Number(form.value.maxCapacity) || 0,
      material: form.value.material,
      propulsion: form.value.propulsion,
      imageUrl: form.value.imageUrl || null
    }
    if (editingId.value) {
      await axios.put(`/api/v1/kayaks/${editingId.value}`, payload)
    } else {
      await axios.post('/api/v1/kayaks/', payload)
    }
    showForm.value = false
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'No se pudo guardar el kayak'
  } finally {
    saving.value = false
  }
}

function askRemove(kayak) {
  confirmingId.value = kayak.id
}

function cancelRemove() {
  confirmingId.value = null
}

async function remove(kayak) {
  try {
    await axios.delete(`/api/v1/kayaks/${kayak.id}`)
    confirmingId.value = null
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'No se pudo eliminar el kayak'
  }
}

onMounted(load)
</script>

<template>
  <div class="kayaks-view">
    <div class="page-head">
      <div>
        <span class="kicker">Equipo</span>
        <h1 class="hemmingway-title">Mis kayaks</h1>
      </div>
      <button class="hemmingway-btn hemmingway-btn-primary" @click="openCreate">
        Agregar kayak
      </button>
    </div>

    <div v-if="loading" class="loading">Cargando kayaks...</div>

    <div v-else-if="error && kayaks.length === 0" class="empty-state">
      <h3>No se pudieron cargar los kayaks</h3>
      <p class="text-muted">{{ error }}</p>
    </div>

    <div v-else-if="kayaks.length === 0" class="empty-state">
      <h3>Todavía no cargaste ningún kayak</h3>
      <p class="text-muted">Cargá el primero para asociarlo a tus salidas y torneos.</p>
    </div>

    <div v-else class="kayak-grid">
      <div v-for="kayak in kayaks" :key="kayak.id" class="hemmingway-card kayak-card">
        <div class="kayak-head">
          <h3>{{ kayak.name }}</h3>
          <span class="hemmingway-badge">{{ kayak.type || 'KAYAK' }}</span>
        </div>

        <dl class="kayak-specs">
          <div v-if="kayak.length">
            <dt>Largo</dt>
            <dd>{{ kayak.length }} cm</dd>
          </div>
          <div v-if="kayak.width">
            <dt>Ancho</dt>
            <dd>{{ kayak.width }} cm</dd>
          </div>
          <div v-if="kayak.weight">
            <dt>Peso</dt>
            <dd>{{ kayak.weight }} kg</dd>
          </div>
          <div v-if="kayak.material">
            <dt>Material</dt>
            <dd>{{ kayak.material }}</dd>
          </div>
          <div v-if="kayak.propulsion">
            <dt>Propulsión</dt>
            <dd>{{ kayak.propulsion }}</dd>
          </div>
        </dl>

        <div class="kayak-actions">
          <button class="hemmingway-btn hemmingway-btn-ghost" @click="openEdit(kayak)">
            Editar
          </button>
          <template v-if="confirmingId === kayak.id">
            <button class="hemmingway-btn hemmingway-btn-secondary" @click="remove(kayak)">
              Confirmar
            </button>
            <button class="hemmingway-btn hemmingway-btn-ghost" @click="cancelRemove">
              Cancelar
            </button>
          </template>
          <button
            v-else
            class="hemmingway-btn hemmingway-btn-secondary"
            @click="askRemove(kayak)"
          >
            Eliminar
          </button>
        </div>
      </div>
    </div>

    <div v-if="showForm" class="modal-backdrop" @click.self="closeForm">
      <form class="hemmingway-card modal" @submit.prevent="save">
        <h2 class="hemmingway-title">{{ editingId ? 'Editar kayak' : 'Agregar kayak' }}</h2>

        <label class="field">
          <span>Nombre</span>
          <input v-model="form.name" class="hemmingway-input" required />
        </label>

        <label class="field">
          <span>Tipo</span>
          <select v-model="form.type" class="hemmingway-input">
            <option v-for="t in TYPES" :key="t" :value="t">{{ t }}</option>
          </select>
        </label>

        <div class="field-row">
          <label class="field">
            <span>Largo (cm)</span>
            <input v-model="form.length" type="number" step="any" class="hemmingway-input" />
          </label>
          <label class="field">
            <span>Ancho (cm)</span>
            <input v-model="form.width" type="number" step="any" class="hemmingway-input" />
          </label>
        </div>

        <div class="field-row">
          <label class="field">
            <span>Peso (kg)</span>
            <input v-model="form.weight" type="number" step="any" class="hemmingway-input" />
          </label>
          <label class="field">
            <span>Capacidad (kg)</span>
            <input v-model="form.maxCapacity" type="number" step="any" class="hemmingway-input" />
          </label>
        </div>

        <label class="field">
          <span>Material</span>
          <select v-model="form.material" class="hemmingway-input">
            <option v-for="m in MATERIALS" :key="m" :value="m">{{ m }}</option>
          </select>
        </label>

        <label class="field">
          <span>Propulsión</span>
          <select v-model="form.propulsion" class="hemmingway-input">
            <option v-for="p in PROPULSIONS" :key="p" :value="p">{{ p }}</option>
          </select>
        </label>

        <p v-if="error" class="form-error">{{ error }}</p>

        <div class="modal-actions">
          <button type="button" class="hemmingway-btn hemmingway-btn-ghost" @click="closeForm">
            Cancelar
          </button>
          <button type="submit" class="hemmingway-btn hemmingway-btn-primary" :disabled="saving">
            {{ saving ? 'Guardando...' : 'Guardar' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}
.kayak-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 18px;
}
.kayak-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
}
.kayak-head h3 {
  margin: 0;
  font-size: 1.1rem;
}
.kayak-specs {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px 16px;
  margin: 0 0 16px;
}
.kayak-specs dt {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--hemmingway-text-muted, #5d7794);
}
.kayak-specs dd {
  margin: 2px 0 0;
  font-family: 'IBM Plex Mono', monospace;
  font-size: 0.9rem;
}
.kayak-actions {
  display: flex;
  gap: 10px;
}
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(3, 8, 15, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  z-index: 50;
}
.modal {
  width: min(520px, 100%);
  max-height: 90vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 14px;
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
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 8px;
}
.form-error {
  color: #ff8a80;
  margin: 0;
  font-size: 0.85rem;
}
.empty-state {
  text-align: center;
  padding: 48px 16px;
}
.empty-state h3 {
  margin: 0 0 6px;
}
.loading {
  padding: 32px 0;
  color: var(--hemmingway-text-muted, #8ba3bd);
}
</style>
