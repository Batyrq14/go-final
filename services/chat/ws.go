package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"qasynda/shared/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for demo
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

type Hub struct {
	clients    map[string]*Client // userID -> Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
	store      *Store
	rmq        *RabbitMQProducer // We will inject this
}

func NewHub(store *Store, rmq *RabbitMQProducer) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte), // Not used directly for P2P but logic handles it
		store:      store,
		rmq:        rmq,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			logger.Info("Client registered: " + client.userID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			logger.Info("Client unregistered: " + client.userID)
		}
	}
}

func (h *Hub) SendPrivateMessage(senderID, receiverID, content string) {
	msg := &Message{
		ID:         uuid.New(),
		SenderID:   uuid.MustParse(senderID),
		ReceiverID: uuid.MustParse(receiverID),
		Content:    content,
		CreatedAt:  time.Now(),
	}

	// 1. Send to RabbitMQ for persistence
	go func() {
		if err := h.rmq.PublishMessage(msg); err != nil {
			logger.Error("failed to publish message", err)
		}
	}()

	// 2. Send to Receiver if online
	h.mu.RLock()
	receiver, ok := h.clients[receiverID]
	h.mu.RUnlock()

	payload, _ := json.Marshal(msg)

	if ok {
		select {
		case receiver.send <- payload:
		default:
			close(receiver.send)
			h.mu.Lock()
			delete(h.clients, receiverID)
			h.mu.Unlock()
		}
	}

	// Also send back to sender for confirmation/UI update if needed,
	// but usually UI handles optimistic update.
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Expecting JSON message: { "receiver_id": "...", "content": "..." }
		var req struct {
			ReceiverID string `json:"receiver_id"`
			Content    string `json:"content"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			logger.Error("invalid message format", err)
			continue
		}

		c.hub.SendPrivateMessage(c.userID, req.ReceiverID, req.Content)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id") // Simple auth for demo
	if userID == "" {
		http.Error(w, "user_id required", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade failed", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), userID: userID}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
