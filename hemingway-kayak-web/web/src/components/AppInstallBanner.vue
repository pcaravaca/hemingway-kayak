<script setup>
import { ref, computed, watch, onMounted } from 'vue'

// En el escritorio convence instalar la app; en el teléfono ya la tiene (o está
// en la tienda). Se avisa una sola vez por navegador.
const DISMISS_KEY = 'hk-app-install-dismissed'

const dismissed = ref(true)

onMounted(() => {
  dismissed.value = localStorage.getItem(DISMISS_KEY) === '1'
})

const isMobile = computed(() => {
  if (typeof navigator === 'undefined') return false
  if (/Android|iPhone|iPad|iPod|Mobile|Opera Mini/i.test(navigator.userAgent)) return true
  const coarse = window.matchMedia && window.matchMedia('(pointer: coarse)').matches
  return coarse && window.innerWidth < 900
})

const visible = computed(() => !isMobile.value && !dismissed.value)

// Mientras está visible, el shell deja el lugar para que la barra no le tape
// el pie a nada.
watch(
  visible,
  (value) => {
    document.body.classList.toggle('has-install-bar', value)
  },
  { immediate: true }
)

function close() {
  dismissed.value = true
  localStorage.setItem(DISMISS_KEY, '1')
}
</script>

<template>
  <div v-if="visible" class="install-bar" role="complementary">
    <div class="install-inner">
      <div class="install-copy">
        <b>La app va mejor en el teléfono.</b>
        <span>
          Tracking con GPS, SOS, spots y capturas con foto. Instalala en Android y dejá
          el navegador para planificar.
        </span>
      </div>
      <div class="install-actions">
        <router-link to="/descargar" class="install-cta">Instalar la app</router-link>
        <button type="button" class="install-close" @click="close" aria-label="Cerrar">×</button>
      </div>
    </div>
  </div>
</template>

<style>
/* Lo pone el componente mientras la barra está visible: así nunca queda nada
   tapado por debajo. */
body.has-install-bar .app-shell {
  padding-bottom: 88px;
}
</style>

<style scoped>
.install-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 60;
  border-top: 1px solid var(--hemmingway-line-2);
  background: rgba(6, 13, 26, .97);
  backdrop-filter: blur(8px);
}

.install-inner {
  max-width: 1240px;
  margin: 0 auto;
  padding: 16px 40px;
  display: flex;
  align-items: center;
  gap: 24px;
}

.install-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13.5px;
  line-height: 1.45;
  color: var(--hemmingway-text-secondary);
}

.install-copy b {
  color: var(--hemmingway-text-primary);
  font-size: 14px;
}

.install-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: none;
  margin-left: auto;
}

.install-cta {
  background: var(--hemmingway-blue);
  color: var(--hemmingway-dark);
  font-size: 12.5px;
  font-weight: 700;
  padding: 10px 16px;
  border-radius: var(--hemmingway-radius-full);
  text-decoration: none;
  white-space: nowrap;
}

.install-cta:hover { box-shadow: var(--hemmingway-shadow-glow); }

.install-close {
  background: transparent;
  border: none;
  color: var(--hemmingway-text-muted);
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  padding: 2px 4px;
}

.install-close:hover { color: var(--hemmingway-text-primary); }

@media (max-width: 700px) {
  .install-inner {
    padding: 14px 20px;
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  .install-actions { margin-left: 0; }
}
</style>
