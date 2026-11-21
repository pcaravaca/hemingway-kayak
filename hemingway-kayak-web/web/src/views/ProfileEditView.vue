<template>
  <div class="edit">
    <router-link to="/profile" class="back">← Volver al perfil</router-link>

    <header class="edit-head">
      <span class="kicker"><i></i>Perfil</span>
      <h1>Editar tu perfil</h1>
      <p class="sub">Cómo te ven los demás kayakers del agua.</p>
    </header>

    <form class="card form" @submit.prevent="handleUpdate">
      <div class="field">
        <span>Foto de perfil</span>
        <div class="avatar-row">
          <img
            v-if="form.avatarUrl"
            :src="form.avatarUrl"
            class="avatar-preview"
            alt="Tu foto de perfil"
          />
          <div v-else class="avatar-preview avatar-empty">{{ inicial }}</div>
          <label
            class="hemmingway-btn hemmingway-btn-secondary avatar-btn"
            :class="{ 'is-disabled': uploading }"
          >
            {{ uploading ? 'Subiendo…' : 'Elegir foto' }}
            <input
              type="file"
              accept="image/*"
              class="avatar-input"
              :disabled="uploading"
              @change="onAvatarSelected"
            />
          </label>
        </div>
        <span v-if="avatarHint" class="field-hint">{{ avatarHint }}</span>
      </div>

      <label class="field">
        <span>Nombre visible</span>
        <input v-model="form.displayName" class="hemmingway-input" placeholder="Tu nombre" />
      </label>

      <label class="field">
        <span>Biografía</span>
        <textarea
          v-model="form.bio"
          class="hemmingway-input"
          rows="4"
          placeholder="Contanos qué pescás y dónde"
        ></textarea>
      </label>

      <label class="switch">
        <input v-model="form.trackingEnabled" type="checkbox" />
        <span class="track"><span class="thumb"></span></span>
        <span class="switch-copy">
          <b>Tracking en vivo</b>
          <em>Permite que otros sigan tu posición en el mapa.</em>
        </span>
      </label>

      <p v-if="message" class="message" :class="{ ok: saved }">{{ message }}</p>

      <div class="actions">
        <button type="submit" class="hemmingway-btn hemmingway-btn-primary" :disabled="saving">
          {{ saving ? 'Guardando…' : 'Guardar cambios' }}
        </button>
        <router-link to="/profile" class="hemmingway-btn hemmingway-btn-secondary">Cancelar</router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const form = ref({ displayName: '', bio: '', trackingEnabled: false, avatarUrl: '' })
const saving = ref(false)
const saved = ref(false)
const message = ref('')
const uploading = ref(false)
const avatarHint = ref('')

// Inicial del nombre, para el placeholder cuando todavía no hay foto.
const inicial = computed(() => {
  const n = form.value.displayName || auth.user?.username || ''
  return (n.trim().charAt(0) || '?').toUpperCase()
})

onMounted(async () => {
  if (auth.user) await auth.fetchUser()
  const u = auth.user || {}
  form.value = {
    displayName: u.displayName || u.username || '',
    bio: u.bio || '',
    trackingEnabled: !!u.trackingEnabled,
    avatarUrl: u.avatarUrl || ''
  }
})

// La foto se sube a /api/v1/uploads (que valida el contenido) y lo que se guarda
// en el perfil es la URL que devuelve; el backend la sirve estática por /uploads.
async function onAvatarSelected(event) {
  const file = event.target.files?.[0]
  if (!file) return
  avatarHint.value = ''
  if (file.size > 8 * 1024 * 1024) {
    avatarHint.value = 'La foto supera los 8 MB.'
    event.target.value = ''
    return
  }

  uploading.value = true
  try {
    const data = new FormData()
    data.append('photo', file)
    const { data: res } = await axios.post('/api/v1/uploads', data)
    form.value.avatarUrl = res.url
    avatarHint.value = 'Foto lista. Guardá los cambios para aplicarla.'
  } catch (e) {
    avatarHint.value = e.response?.data?.error || 'No se pudo subir la foto.'
    event.target.value = ''
  } finally {
    uploading.value = false
  }
}

async function handleUpdate() {
  saving.value = true
  message.value = ''
  const result = await auth.updateProfile({
    displayName: form.value.displayName,
    bio: form.value.bio,
    trackingEnabled: form.value.trackingEnabled,
    avatarUrl: form.value.avatarUrl || null
  })
  saving.value = false

  if (result.success) {
    saved.value = true
    message.value = 'Perfil guardado.'
    setTimeout(() => router.push('/profile'), 700)
  } else {
    saved.value = false
    message.value = result.error || auth.error || 'No se pudo guardar.'
  }
}
</script>

<style scoped>
.edit { padding: 40px 0 24px; max-width: 640px; }

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

.edit-head { padding-bottom: 24px; }

.edit-head h1 {
  font-family: var(--hemmingway-font-display);
  font-weight: 300;
  font-size: 36px;
  letter-spacing: -.02em;
  margin: 12px 0 8px;
}

.sub { color: var(--hemmingway-text-secondary); font-size: 15px; }

.form { display: flex; flex-direction: column; gap: 20px; }

.field { display: flex; flex-direction: column; gap: 8px; }

.field span {
  font-family: var(--hemmingway-font-mono);
  font-size: 11px;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: var(--hemmingway-text-muted);
}

.field-hint {
  font-family: var(--hemmingway-font-mono);
  font-size: 12px;
  color: var(--hemmingway-text-secondary);
}

.avatar-row { display: flex; align-items: center; gap: 16px; }

.avatar-preview {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--hemmingway-line-2);
  background: var(--hemmingway-bg-tertiary);
}

.avatar-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--hemmingway-font-display);
  font-size: 28px;
  color: var(--hemmingway-blue);
}

/* El input va escondido pero accesible: el label es el botón visible. */
.avatar-input { position: absolute; opacity: 0; width: 0; height: 0; }

.avatar-btn { position: relative; cursor: pointer; }

.avatar-btn.is-disabled { opacity: .6; pointer-events: none; }

textarea.hemmingway-input { resize: vertical; min-height: 96px; }

.switch { display: flex; align-items: center; gap: 14px; cursor: pointer; }

.switch input { position: absolute; opacity: 0; width: 0; height: 0; }

.switch .track {
  flex: none;
  width: 44px;
  height: 24px;
  border-radius: var(--hemmingway-radius-full);
  border: 1px solid var(--hemmingway-line-2);
  background: var(--hemmingway-bg-tertiary);
  position: relative;
  transition: all var(--hemmingway-transition-fast);
}

.switch .thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--hemmingway-text-muted);
  transition: all var(--hemmingway-transition-fast);
}

.switch input:checked + .track { background: rgba(56,189,248,.2); border-color: var(--hemmingway-blue); }
.switch input:checked + .track .thumb { left: 22px; background: var(--hemmingway-blue); }

.switch-copy b { display: block; font-size: 14.5px; }
.switch-copy em {
  font-style: normal;
  font-size: 12.5px;
  color: var(--hemmingway-text-secondary);
}

.message { font-size: 14px; color: var(--hemmingway-accent-error); margin: 0; }
.message.ok { color: var(--hemmingway-accent-success); }

.actions { display: flex; gap: 12px; flex-wrap: wrap; }
</style>
