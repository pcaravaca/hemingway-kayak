package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type string      `json:"type"`
	Room string      `json:"room,omitempty"`
	Data interface{} `json:"data"`
}

type roomEvent struct {
	client *websocket.Conn
	room   string
}

// Hub administra clientes y salas. Todas las estructuras se tocan sólo desde
// el goroutine de Run, o bajo el mutex, para no correr carreras.
type Hub struct {
	mu        sync.RWMutex
	clients   map[*websocket.Conn]map[string]bool
	broadcast chan Message
	register  chan *websocket.Conn
	unregister chan *websocket.Conn
	join      chan roomEvent
	leave     chan roomEvent
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]map[string]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		join:       make(chan roomEvent),
		leave:      make(chan roomEvent),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = make(map[string]bool)
			h.mu.Unlock()
			log.Printf("Client connected. Total: %d", h.clientCount())

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", h.clientCount())

		case ev := <-h.join:
			h.mu.Lock()
			if rooms, ok := h.clients[ev.client]; ok {
				rooms[ev.room] = true
			}
			h.mu.Unlock()
			log.Printf("Client joined room %s", ev.room)

		case ev := <-h.leave:
			h.mu.Lock()
			if rooms, ok := h.clients[ev.client]; ok {
				delete(rooms, ev.room)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}
			h.broadcastTo(message.Room, data)
		}
	}
}

func (h *Hub) clientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) broadcastTo(room string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client, rooms := range h.clients {
		if room != "" && !rooms[room] {
			continue
		}
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error sending message: %v", err)
			// No se puede cerrar/borrar aquí con el RLock tomado: se deja para
			// el unregister que produce el cierre real de la conexión.
			_ = err
		}
	}
}

func (h *Hub) HandleWebSocket(c *websocket.Conn) {
	h.register <- c
	defer func() { h.unregister <- c }()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		var incoming struct {
			Action string `json:"action"`
			Room   string `json:"room"`
		}
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		switch incoming.Action {
		case "subscribe":
			if incoming.Room != "" {
				h.join <- roomEvent{client: c, room: incoming.Room}
			}
		case "unsubscribe":
			if incoming.Room != "" {
				h.leave <- roomEvent{client: c, room: incoming.Room}
			}
		}
	}
}

// Broadcast envía un mensaje a una sala (o a todos si room está vacío).
func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

// BroadcastTo emite un evento de tracking a la sala de un torneo/usuario.
func (h *Hub) BroadcastTo(room string, eventType string, data interface{}) {
	h.Broadcast(Message{Type: eventType, Room: room, Data: data})
}

// UpgradeMiddleware existe por compatibilidad; el middleware real está en
// internal/api/server.go (websocket.New).
func UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("allowed", true)
		return c.Next()
	}
}
