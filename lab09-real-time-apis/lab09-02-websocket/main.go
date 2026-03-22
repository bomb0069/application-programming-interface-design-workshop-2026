package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan Message
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	history    []Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		history:    make([]Message, 0),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// Send chat history to new client
			for _, msg := range h.history {
				client.send <- msg
			}

			// Broadcast join message
			h.broadcast <- Message{
				Type:      "system",
				Content:   client.username + " joined the chat",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

			h.broadcast <- Message{
				Type:      "system",
				Content:   client.username + " left the chat",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}

		case msg := <-h.broadcast:
			h.mu.Lock()
			if msg.Type == "message" {
				h.history = append(h.history, msg)
				if len(h.history) > 100 {
					h.history = h.history[1:]
				}
			}
			h.mu.Unlock()

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

var hub = NewHub()

func main() {
	go hub.Run()

	r := chi.NewRouter()

	r.Get("/ws", handleWebSocket)
	r.Get("/stats", handleStats)
	r.Handle("/*", http.FileServer(http.Dir("static")))

	log.Println("Server starting on :8080")
	log.Println("Open http://localhost:8080 in your browser")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		conn:     conn,
		username: username,
		send:     make(chan Message, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var input struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(rawMsg, &input); err != nil || input.Content == "" {
			continue
		}

		hub.broadcast <- Message{
			Type:      "message",
			Username:  c.username,
			Content:   input.Content,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	usernames := []string{}
	for client := range hub.clients {
		usernames = append(usernames, client.username)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"online_users":  len(hub.clients),
		"usernames":     usernames,
		"message_count": len(hub.history),
	})
}
