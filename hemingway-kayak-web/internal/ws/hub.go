package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Error sending message: %v", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) HandleWebSocket(c *websocket.Conn) {
	h.register <- c

	defer func() {
		h.unregister <- c
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var incoming map[string]interface{}
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		action, _ := incoming["action"].(string)
		switch action {
		case "subscribe":
			log.Printf("Client subscribed to tracking updates")
		case "unsubscribe":
			log.Printf("Client unsubscribed from tracking updates")
		}
	}
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

func UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("allowed", true)
		return c.Next()
	}
}
