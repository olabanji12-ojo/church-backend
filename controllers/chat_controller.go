package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/olabanji12-ojo/church-backend/hub"
	"github.com/olabanji12-ojo/church-backend/middleware"
	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/olabanji12-ojo/church-backend/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow any origin for now (CORS is handled globally, but WS needs specific checks if strict)
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatController struct {
	Hub            *hub.Hub
	MessageService *services.MessageService
}

func NewChatController(h *hub.Hub, messageService *services.MessageService) *ChatController {
	return &ChatController{
		Hub:            h,
		MessageService: messageService,
	}
}

// ServeWS handles WebSocket requests from the clients.
func (cc *ChatController) ServeWS(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate via token in query string (ws://...?auth_token=JWT)
	tokenString := r.URL.Query().Get("auth_token")
	if tokenString == "" {
		http.Error(w, "Missing auth_token", http.StatusUnauthorized)
		return
	}

	userID, err := utils.ValidateJWT(tokenString)
	if err != nil || userID == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 2. Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Failed to upgrade to websocket: %v", err)
		return
	}

	// 3. Register client in Hub
	client := hub.NewClient(cc.Hub, userID, conn, cc.MessageService)
	cc.Hub.Register(client)

	// 4. Start concurrent pumps
	go client.WritePump()
	go client.ReadPump()
}

// GetMessagesHandler fetches the historical messages for a specific match or target user
func (cc *ChatController) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	senderOID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	targetIDStr := mux.Vars(r)["target_user_id"]
	targetOID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	// 1. Try finding messages by target user ID first
	messages, err := cc.MessageService.GetChatHistoryByTarget(senderOID, targetOID)
	if err != nil || len(messages) == 0 {
		// 2. Fallback: Check if targetIDStr was actually a Match ID directly
		matchMessages, err2 := cc.MessageService.GetChatHistory(targetOID)
		if err2 == nil && len(matchMessages) > 0 {
			messages = matchMessages
		}
	}

	if messages == nil {
		messages = []models.Message{}
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": messages,
	})
}
