<template>
  <div class="p-8">
    <h1 class="text-2xl font-bold">Detalle del Viaje: {{ $route.params.id }}</h1>
    <div class="mt-4">
      <p>Cargando detalles del viaje...</p>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const route = useRoute()
const details = ref(null)
onMounted(async () => {
  // fetch details from API
  details.value = { title: 'Ejemplo de Viaje', description: 'Cargando...' }
})
</script>
