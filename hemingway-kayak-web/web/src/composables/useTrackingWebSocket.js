import { ref, onUnmounted } from 'vue'

export function useTrackingWebSocket() {
  const connected = ref(false)
  const trackers = ref([])
  const tournamentUpdates = ref(null)
  let ws = null

  function connect() {
    // En produccion el WS queda en el mismo origin que la SPA (nginx hace proxy);
    // VITE_WS_URL sigue permitiendo fijarlo a mano.
    const wsUrl = import.meta.env.VITE_WS_URL
      || `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws/tracking`
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      connected.value = true
      console.log('Tracking WebSocket connected')
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
      case 'tournament':
        tournamentUpdates.value = message.data
        break
    }
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
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: 'subscribe', tournamentId }))
    }
  }

  function unsubscribeTournament(tournamentId) {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: 'unsubscribe', tournamentId }))
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