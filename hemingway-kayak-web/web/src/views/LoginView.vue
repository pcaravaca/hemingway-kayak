<template>
  <div class="login-view">
    <div class="login-card hemingway-card">
      <div class="login-header">
        <h1 class="hemmingway-title">Bienvenido de nuevo</h1>
        <p class="login-subtitle">Entrá para seguir tus salidas</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email" class="form-label">Correo</label>
          <input
            id="email"
            v-model="email"
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
            v-model="password"
            type="password"
            class="hemmingway-input"
            placeholder="Tu contraseña"
            required
          />
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <button
          type="submit"
          class="hemmingway-btn hemmingway-btn-primary hemmingway-btn-full"
          :disabled="loading"
        >
          {{ loading ? 'Entrando...' : 'Entrar' }}
        </button>
      </form>

      <div class="login-footer">
        <p class="text-muted">
          ¿No tenés cuenta?
          <router-link to="/register" class="link">Creala</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref(null)

async function handleLogin() {
  loading.value = true
  error.value = null

  // authStore.login() devuelve { success, error }: sin el .success este if
  // era siempre verdadero y el login fallido también redirigía al home.
  const result = await authStore.login(email.value, password.value)

  if (result.success) {
    router.push('/')
  } else {
    error.value = result.error || authStore.error || 'No pudimos ingresar'
  }

  loading.value = false
}
</script>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
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

.hemingway-btn-full {
  width: 100%;
  margin-top: var(--hemmingway-space-md);
}

.login-footer {
  text-align: center;
  margin-top: var(--hemmingway-space-xl);
}

.link {
  color: var(--hemmingway-neon);
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}
</style>