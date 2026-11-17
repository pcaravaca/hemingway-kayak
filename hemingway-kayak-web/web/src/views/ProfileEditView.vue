<template>
  <div class="edit">
    <router-link to="/profile" class="back">← Volver al perfil</router-link>

    <header class="edit-head">
      <span class="kicker"><i></i>Perfil</span>
      <h1>Editar tu perfil</h1>
      <p class="sub">Cómo te ven los demás kayakers del agua.</p>
    </header>

    <form class="card form" @submit.prevent="handleUpdate">
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
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const form = ref({ displayName: '', bio: '', trackingEnabled: false })
const saving = ref(false)
const saved = ref(false)
const message = ref('')

onMounted(async () => {
  if (auth.user) await auth.fetchUser()
  const u = auth.user || {}
  form.value = {
    displayName: u.displayName || u.username || '',
    bio: u.bio || '',
    trackingEnabled: !!u.trackingEnabled
  }
})

async function handleUpdate() {
  saving.value = true
  message.value = ''
  const result = await auth.updateProfile({
    displayName: form.value.displayName,
    bio: form.value.bio,
    trackingEnabled: form.value.trackingEnabled
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
