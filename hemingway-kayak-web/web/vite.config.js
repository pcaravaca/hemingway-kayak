import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// El backend puede correr en otro host/puerto; se sobrescribe con VITE_API_TARGET.
const apiTarget = process.env.VITE_API_TARGET || 'http://localhost:8080'
const wsTarget = process.env.VITE_WS_TARGET || apiTarget.replace(/^http/, 'ws')

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    host: '0.0.0.0',
    // NTFS: los mtimes son poco fiables y el watcher puede no invalidar.
    watch: {
      usePolling: true,
      interval: 300,
    },
    proxy: {
      '/api': {
        target: apiTarget,
        changeOrigin: true,
      },
      '/ws': {
        target: wsTarget,
        ws: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
  },
})
