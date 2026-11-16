<template>
  <div class="p-8 max-w-md mx-auto">
    <h1 class="text-2xl font-bold mb-4">Editar Perfil</h1>
    <form @submit.prevent="handleUpdate" class="space-y-4">
      <input v-model="form.displayName" placeholder="Nombre de usuario" class="w-full p-2 border rounded" />
      <textarea v-model="form.bio" placeholder="Biografía" class="w-full p-2 border rounded"></textarea>
      <button type="submit" class="bg-blue-600 text-white p-2 rounded">Guardar</button>
    </form>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
const auth = useAuthStore()
const router = useRouter()
const form = ref({ displayName: '', bio: '' })
onMounted(async () => {
  if (auth.user) await auth.fetchUser()
  form.value = auth.user || { displayName: '', bio: '' }
})
async function handleUpdate() {
  await auth.updateProfile(form.value)
  router.push('/profile')
}
</script>
