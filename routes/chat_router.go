package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/middleware"
)

func ChatRoutes(router *mux.Router, chatController *controllers.ChatController) {
	// Notice this is not /api/v1, but /ws/v1 because it is a WebSocket route
	ws := router.PathPrefix("/ws/v1/chat").Subrouter()
	ws.HandleFunc("", chatController.ServeWS).Methods("GET")

	// REST endpoint for fetching chat history
	api := router.PathPrefix("/api/v1/chat").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/{target_user_id}/messages", chatController.GetMessagesHandler).Methods("GET")
}
