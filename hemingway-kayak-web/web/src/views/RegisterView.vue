<template>
  <div class="p-8 max-w-md mx-auto">
    <h1 class="text-2xl font-bold mb-4">Registro</h1>
    <form @submit.prevent="handleRegister" class="space-y-4">
      <input v-model="form.name" placeholder="Nombre" class="w-full p-2 border rounded" required />
      <input v-model="form.email" type="email" placeholder="Email" class="w-full p-2 border rounded" required />
      <input v-model="form.password" type="password" placeholder="Contraseña" class="w-full p-2 border rounded" required />
      <button type="submit" class="bg-blue-600 text-white p-2 rounded">Registrarse</button>
    </form>
    <p class="mt-4">¿Ya tienes cuenta? <router-link to="/login">Inicia sesión</router-link></p>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
const auth = useAuthStore()
const router = useRouter()
const form = ref({ name: '', email: '', password: '' })
async function handleRegister() {
  const { success, error } = await auth.register(form.value)
  if (success) router.push('/login')
  else alert(error)
}
</script>
