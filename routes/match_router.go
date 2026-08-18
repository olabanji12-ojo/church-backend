package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/middleware"
)

func MatchRoutes(router *mux.Router, matchController *controllers.MatchController) {
	match := router.PathPrefix("/api/v1/matches").Subrouter()
	
	match.Use(middleware.AuthMiddleware)

	match.HandleFunc("/requests", matchController.GetPendingLikesHandler).Methods("GET")
	match.HandleFunc("/{target_user_id}/like", matchController.SwipeRightHandler).Methods("POST")
	match.HandleFunc("/{target_user_id}/pass", matchController.SwipeLeftHandler).Methods("POST")
	match.HandleFunc("/{target_user_id}", matchController.UnmatchHandler).Methods("DELETE")
	match.HandleFunc("", matchController.GetMatchesHandler).Methods("GET")
}
