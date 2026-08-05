package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
)

func AuthRoutes(router *mux.Router, authController *controllers.AuthController) {
	auth := router.PathPrefix("/api/v1/auth").Subrouter()

	auth.HandleFunc("/register", authController.RegisterHandler).Methods("POST")
	auth.HandleFunc("/login", authController.LoginHandler).Methods("POST")
	auth.HandleFunc("/logout", authController.LogoutHandler).Methods("POST")
	auth.HandleFunc("/social", authController.SocialAuthHandler).Methods("POST")
	auth.HandleFunc("/guest", authController.GuestLoginHandler).Methods("POST")
}
