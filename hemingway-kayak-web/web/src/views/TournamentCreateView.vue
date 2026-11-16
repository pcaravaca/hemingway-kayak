<template>
  <div class="p-8 max-w-md mx-auto">
    <h1 class="text-2xl font-bold mb-4">Crear Torneo</h1>
    <form @submit.prevent="handleCreate" class="space-y-4">
      <input v-model="form.name" placeholder="Nombre del Torneo" class="w-full p-2 border rounded" required />
      <input v-model="form.type" placeholder="Tipo (Kayak, etc)" class="w-full p-2 border rounded" required />
      <button type="submit" class="bg-blue-600 text-white p-2 rounded">Crear</button>
    </form>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const form = ref({ name: '', type: '' })
async function handleCreate() {
  console.log('Creating:', form.value)
  router.push('/tournaments')
}
</script>
