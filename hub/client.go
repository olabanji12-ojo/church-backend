package hub

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub            *Hub
	UserID         string
	conn           *websocket.Conn
	send           chan []byte
	MessageService *services.MessageService
}

func NewClient(hub *Hub, userID string, conn *websocket.Conn, messageService *services.MessageService) *Client {
	return &Client{
		hub:            hub,
		UserID:         userID,
		conn:           conn,
		send:           make(chan []byte, 256),
		MessageService: messageService,
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("WebSocket read error: %v", err)
			}
			break
		}

		var payload MessagePayload
		if err := json.Unmarshal(message, &payload); err != nil {
			logrus.Warnf("Invalid socket message format: %v", err)
			continue
		}

		// 1. Validate Target User & Save to Database
		senderOID, _ := primitive.ObjectIDFromHex(c.UserID)
		targetOID, _ := primitive.ObjectIDFromHex(payload.TargetUserID)
		
		newMsg, err := c.MessageService.SendMessageToTargetUser(senderOID, targetOID, payload.Content, payload.Type)
		if err != nil {
			logrus.Errorf("Failed to save message to DB or invalid target: %v", err)
			continue
		}

		// 2. Wrap into RedisMessage to include the recipient ID
		redisMsg := RedisMessage{
			Message:     *newMsg,
			RecipientID: payload.TargetUserID, // The recipient is the target user
		}

		// 3. Broadcast via Redis
		c.hub.PublishToRedis(redisMsg)
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
