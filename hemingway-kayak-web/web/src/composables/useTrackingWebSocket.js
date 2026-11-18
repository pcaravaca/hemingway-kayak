import { ref, onUnmounted } from 'vue'

export function useTrackingWebSocket() {
  const connected = ref(false)
  const trackers = ref([])
  const tournamentUpdates = ref(null)
  let ws = null
  let pendingRoom = null

  function connect() {
    // En produccion el WS queda en el mismo origin que la SPA (nginx hace proxy);
    // VITE_WS_URL sigue permitiendo fijarlo a mano.
    const wsUrl = import.meta.env.VITE_WS_URL
      || `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws/tracking`
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      connected.value = true
      console.log('Tracking WebSocket connected')
      // La suscripción puede pedirse antes de que el socket exista (el componente
      // se monta, pide el torneo y recién después llama a connect()); sin esto,
      // el primer subscribe se perdía y la sala quedaba vacía.
      if (pendingRoom) {
        ws.send(JSON.stringify({ action: 'subscribe', room: pendingRoom }))
      }
    }

    ws.onclose = () => {
      connected.value = false
      console.log('Tracking WebSocket disconnected')
      setTimeout(connect, 5000)
    }

    ws.onerror = (err) => {
      console.error('WebSocket error:', err)
    }

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        handleMessage(message)
      } catch (e) {
        console.error('Failed to parse message:', e)
      }
    }
  }

  function handleMessage(message) {
    switch (message.type) {
      case 'tracking':
        updateTracker(message.data)
        break
      case 'tracking-stop':
        removeTracker(message.data?.userId)
        break
      case 'tournament':
        tournamentUpdates.value = message.data
        break
    }
  }

  function removeTracker(userId) {
    if (!userId) return
    trackers.value = trackers.value.filter(t => t.userId !== userId)
  }

  function updateTracker(data) {
    const index = trackers.value.findIndex(t => t.userId === data.userId)
    if (index >= 0) {
      trackers.value[index] = { ...trackers.value[index], ...data, lastUpdate: Date.now() }
    } else {
      trackers.value.push({ ...data, lastUpdate: Date.now() })
    }
  }

  function subscribeTournament(tournamentId) {
    // El hub espera `room` (internal/ws/hub.go); antes se mandaba `tournamentId`
    // y la sala quedaba vacía, así que la suscripción nunca tenía efecto.
    pendingRoom = tournamentId
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: 'subscribe', room: tournamentId }))
    }
  }

  function unsubscribeTournament(tournamentId) {
    if (pendingRoom === tournamentId) {
      pendingRoom = null
    }
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: 'unsubscribe', room: tournamentId }))
    }
  }

  function disconnect() {
    if (ws) {
      ws.close()
      ws = null
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    trackers,
    tournamentUpdates,
    connect,
    disconnect,
    subscribeTournament,
    unsubscribeTournament
  }
}