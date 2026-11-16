<template>
  <div class="login-view">
    <div class="login-card hemmingway-card">
      <div class="login-header">
        <h1 class="hemmingway-title">Creá tu cuenta</h1>
        <p class="login-subtitle">Empezá a registrar tus salidas</p>
      </div>

      <form @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label for="name" class="form-label">Nombre</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            class="hemmingway-input"
            placeholder="Tu nombre"
            required
          />
        </div>

        <div class="form-group">
          <label for="email" class="form-label">Correo</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            class="hemmingway-input"
            placeholder="tu@correo.com"
            required
          />
        </div>

        <div class="form-group">
          <label for="password" class="form-label">Contraseña</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            class="hemmingway-input"
            placeholder="Mínimo 8 caracteres"
            minlength="8"
            required
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="ok" class="ok-message">{{ ok }}</div>

        <button
          type="submit"
          class="hemmingway-btn hemmingway-btn-primary hemmingway-btn-full"
          :disabled="loading"
        >
          {{ loading ? 'Creando...' : 'Crear cuenta' }}
        </button>
      </form>

      <div class="login-footer">
        <p class="text-muted">
          ¿Ya tenés cuenta?
          <router-link to="/login" class="link">Entrá</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const form = ref({ name: '', email: '', password: '' })
const loading = ref(false)
const error = ref(null)
const ok = ref(null)

async function handleRegister() {
  loading.value = true
  error.value = null
  ok.value = null

  const result = await auth.register({
    displayName: form.value.name,
    email: form.value.email,
    username: form.value.email.split('@')[0],
    password: form.value.password
  })

  if (result.success) {
    ok.value = result.message || 'Cuenta creada. Ahora podés entrar.'
    setTimeout(() => router.push('/login'), 1200)
  } else {
    error.value = result.error || 'No pudimos crear la cuenta'
  }

  loading.value = false
}
</script>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: calc(100vh - 220px);
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: var(--hemmingway-space-2xl);
}

.login-header {
  text-align: center;
  margin-bottom: var(--hemmingway-space-xl);
}

.login-subtitle {
  color: var(--hemmingway-text-secondary);
  margin-top: var(--hemmingway-space-sm);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--hemmingway-space-sm);
}

.form-label {
  font-size: var(--hemmingway-font-size-sm);
  font-weight: 600;
  color: var(--hemmingway-text-secondary);
}

.error-message {
  padding: var(--hemmingway-space-md);
  background: rgba(255, 107, 107, 0.1);
  border: 1px solid var(--hemmingway-accent-error);
  border-radius: var(--hemmingway-radius-md);
  color: var(--hemmingway-accent-error);
  font-size: var(--hemmingway-font-size-sm);
}

.ok-message {
  padding: var(--hemmingway-space-md);
  background: rgba(52, 211, 153, 0.1);
  border: 1px solid var(--hemmingway-accent-success);
  border-radius: var(--hemmingway-radius-md);
  color: var(--hemmingway-accent-success);
  font-size: var(--hemmingway-font-size-sm);
}

.hemmingway-btn-full {
  width: 100%;
  margin-top: var(--hemmingway-space-md);
}

.login-footer {
  text-align: center;
  margin-top: var(--hemmingway-space-xl);
}

.link {
  color: var(--hemmingway-blue);
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}
</style>
