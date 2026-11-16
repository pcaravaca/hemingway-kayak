<template>
  <div class="download-view">
    <div class="download-card hemingway-card">
      <div class="download-header">
        <img src="/logo.svg" alt="Hemingway Kayak" class="download-logo" />
        <h1 class="hemingway-title">Descargá la app</h1>
        <p class="download-subtitle">
          Llevá tus viajes, waypoints y torneos en el bolsillo.
        </p>
      </div>

      <a
        class="hemingway-btn hemingway-btn-primary download-btn"
        href="/HemingwayKayak.apk"
        download="HemingwayKayak.apk"
      >
        Descargar para Android
      </a>

      <p class="download-meta text-muted">
        APK · {{ sizeLabel }} · Android 7.0 o superior
      </p>

      <ol class="steps">
        <li>Abrí el archivo descargado desde la notificación o tu carpeta de Descargas.</li>
        <li>Si el teléfono lo pide, permití la instalación desde este navegador.</li>
        <li>Abrí la app y entrá con tu cuenta de Hemingway Kayak.</li>
      </ol>

      <p class="download-note text-muted">
        ¿Ya tenés la app? <router-link to="/login" class="link">Iniciá sesión</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const sizeLabel = ref('—')

onMounted(async () => {
  try {
    const res = await fetch('/HemingwayKayak.apk', { method: 'HEAD' })
    const len = Number(res.headers.get('content-length') || 0)
    if (len) {
      sizeLabel.value = (len / (1024 * 1024)).toFixed(1) + ' MB'
    }
  } catch {
    sizeLabel.value = 'APK'
  }
})
</script>

<style scoped>
.download-view {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: calc(100vh - 200px);
}

.download-card {
  width: 100%;
  max-width: 460px;
  padding: var(--hemmingway-space-2xl);
  text-align: center;
}

.download-header {
  margin-bottom: var(--hemmingway-space-xl);
}

.download-logo {
  width: 84px;
  height: 84px;
  margin-bottom: var(--hemmingway-space-md);
}

.download-subtitle {
  color: var(--hemmingway-text-secondary);
  margin-top: var(--hemmingway-space-sm);
}

.download-btn {
  display: block;
  width: 100%;
  text-decoration: none;
  font-size: var(--hemmingway-font-size-lg);
  padding: 14px 20px;
  background: linear-gradient(135deg, #4DFFA9 0%, #40E0D0 100%);
  color: #1A1B1E;
  font-weight: 700;
  border-radius: var(--hemmingway-radius-lg);
}

.download-meta {
  margin-top: var(--hemmingway-space-md);
  font-size: var(--hemmingway-font-size-sm);
}

.steps {
  text-align: left;
  margin: var(--hemmingway-space-xl) 0 0;
  padding-left: var(--hemmingway-space-xl);
  color: var(--hemingway-text-secondary);
  font-size: var(--hemmingway-font-size-sm);
  line-height: 1.7;
}

.download-note {
  margin-top: var(--hemmingway-space-xl);
  font-size: var(--hemmingway-font-size-sm);
}

.link {
  color: var(--hemmingway-neon);
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}
</style>
