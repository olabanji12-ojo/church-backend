package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/middleware"
)

func PrayerRoutes(router *mux.Router, prayerController *controllers.PrayerController) {
	prayer := router.PathPrefix("/api/v1/prayers").Subrouter()
	
	prayer.Use(middleware.AuthMiddleware)

	prayer.HandleFunc("", prayerController.GetPrayersHandler).Methods("GET")
	prayer.HandleFunc("", prayerController.PostPrayerHandler).Methods("POST")
	prayer.HandleFunc("/{id}/amen", prayerController.AddAmenHandler).Methods("POST")
}
