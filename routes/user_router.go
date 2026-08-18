package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/middleware"
)

func UserRoutes(router *mux.Router, userController *controllers.UserController) {
	user := router.PathPrefix("/api/v1/users").Subrouter()
	
	// Apply Auth Middleware to all user routes
	user.Use(middleware.AuthMiddleware)

	user.HandleFunc("/me", userController.GetCurrentUserHandler).Methods("GET")
	user.HandleFunc("/me", userController.UpdateCurrentUserHandler).Methods("PATCH")
	user.HandleFunc("/me/faith-profile", userController.UpdateFaithProfileHandler).Methods("PATCH")
	user.HandleFunc("/push-token", userController.SavePushToken).Methods("POST")

	user.HandleFunc("/scenarios/questions", userController.GetScenarioQuestionsHandler).Methods("GET")
	user.HandleFunc("/scenarios/answer", userController.AnswerScenarioHandler).Methods("POST")
}

// DiscoverRoutes (since it's slightly separate conceptually but uses UserController)
func DiscoverRoutes(router *mux.Router, userController *controllers.UserController) {
	discover := router.PathPrefix("/api/v1/discover").Subrouter()
	
	discover.Use(middleware.AuthMiddleware)

	discover.HandleFunc("", userController.GetDiscoveryFeedHandler).Methods("GET")
}
