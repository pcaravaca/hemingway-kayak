import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

axios.defaults.baseURL = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const isAuthenticated = computed(() => !!user.value)
  const token = ref(localStorage.getItem('token') || null)
  const loading = ref(false)
  const error = ref(null)

  axios.defaults.headers.common['Authorization'] = token.value ? `Bearer ${token.value}` : ''

  async function fetchUser() {
    if (!token.value) return
    try {
      const response = await axios.get('/api/v1/users/me')
      user.value = response.data
    } catch (e) {
      logout()
    }
  }

  async function login(email, password) {
    loading.value = true
    error.value = null
    try {
      const response = await axios.post('/api/v1/auth/login', { email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', token.value)
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      return { success: true }
    } catch (e) {
      error.value = e.response?.data?.error || e.response?.data?.message || 'No pudimos ingresar'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function register(userData) {
    loading.value = true
    error.value = null
    try {
      await axios.post('/api/v1/auth/register', userData)
      return { success: true, message: 'Cuenta creada. Ahora podés entrar.' }
    } catch (e) {
      error.value = e.response?.data?.error || e.response?.data?.message || 'No pudimos crear la cuenta'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function updateProfile(userData) {
    loading.value = true
    error.value = null
    try {
      const id = user.value?.id
      if (!id) throw new Error('Not authenticated')
      await axios.put(`/api/v1/users/${id}`, userData)
      await fetchUser()
      return { success: true }
    } catch (e) {
      error.value = e.response?.data?.error || e.response?.data?.message || 'Failed to update profile'
      return { success: false }
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await axios.post('/api/v1/auth/logout').catch(() => {})
    } finally {
      user.value = null
      token.value = null
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
    }
  }

  return { user, isAuthenticated, token, loading, error, login, register, updateProfile, logout, fetchUser }
})
