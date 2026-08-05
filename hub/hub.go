package hub

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/olabanji12-ojo/church-backend/database"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/sirupsen/logrus"
)

// MessagePayload represents the incoming JSON from the WebSocket client
type MessagePayload struct {
	TargetUserID string `json:"target_user_id"` // Used to be match_id
	Content      string `json:"content"`
	Type         string `json:"type,omitempty"`
}

// RedisMessage represents the complete message sent globally over Redis
type RedisMessage struct {
	Message     models.Message `json:"message"`
	RecipientID string         `json:"recipient_id"`
}

// Hub maintains the set of active clients and handles Redis Pub/Sub
type Hub struct {
	// Registered clients mapped by UserID. Value is a map in case a user has multiple devices connected
	clients map[string]map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex to protect clients map
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register adds a new client to the hub
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

func (h *Hub) Run() {
	// Start listening to the global Redis chat channel
	go h.subscribeToRedis()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			logrus.Infof("Client registered: %s", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if connections, ok := h.clients[client.UserID]; ok {
				if _, ok := connections[client]; ok {
					delete(connections, client)
					close(client.send)
					if len(connections) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			logrus.Infof("Client unregistered: %s", client.UserID)
		}
	}
}

// PublishToRedis is called by a client when they send a message
func (h *Hub) PublishToRedis(payload RedisMessage) {
	data, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Error marshaling payload for redis: %v", err)
		return
	}

	// Publish to a global chat channel.
	err = database.RedisClient.Publish(context.Background(), "global:chat", data).Err()
	if err != nil {
		logrus.Errorf("Error publishing to redis: %v", err)
	}
}

// subscribeToRedis listens for incoming messages from ANY server instance
func (h *Hub) subscribeToRedis() {
	pubsub := database.RedisClient.Subscribe(context.Background(), "global:chat")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var redisMsg RedisMessage
		if err := json.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
			logrus.Errorf("Error unmarshaling redis message: %v", err)
			continue
		}

		// We send the FULL RedisMessage to the frontend, so it knows the RecipientID!
		messageBytes := []byte(msg.Payload)

		h.mu.RLock()
		if connections, ok := h.clients[redisMsg.RecipientID]; ok {
			for client := range connections {
				select {
				case client.send <- messageBytes:
				default:
					close(client.send)
					delete(connections, client)
				}
			}
		}
		// Also send it back to the Sender (if they have multiple devices open)
		if connections, ok := h.clients[redisMsg.Message.SenderID.Hex()]; ok {
			for client := range connections {
				select {
				case client.send <- messageBytes:
				default:
					close(client.send)
					delete(connections, client)
				}
			}
		}
		h.mu.RUnlock()
	}
}
