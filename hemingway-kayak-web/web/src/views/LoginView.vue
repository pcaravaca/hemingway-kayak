<template>
  <div class="login-view">
    <div class="login-card hemingway-card">
      <div class="login-header">
        <h1 class="hemingway-title">Welcome Back</h1>
        <p class="login-subtitle">Sign in to continue your adventures</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email" class="form-label">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            class="hemingway-input"
            placeholder="your@email.com"
            required
          />
        </div>

        <div class="form-group">
          <label for="password" class="form-label">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="hemingway-input"
            placeholder="Your password"
            required
          />
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <button
          type="submit"
          class="hemingway-btn hemingway-btn-primary hemingway-btn-full"
          :disabled="loading"
        >
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>

      <div class="login-footer">
        <p class="text-muted">
          Don't have an account?
          <router-link to="/register" class="link">Create one</router-link>
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

  const success = await authStore.login(email.value, password.value)

  if (success) {
    router.push('/')
  } else {
    error.value = authStore.error || 'Login failed'
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