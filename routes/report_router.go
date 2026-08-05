package routes

import (
	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/controllers"
	"github.com/olabanji12-ojo/church-backend/middleware"
)

func ReportRoutes(router *mux.Router, reportController *controllers.ReportController) {
	report := router.PathPrefix("/api/v1/reports").Subrouter()
	
	report.Use(middleware.AuthMiddleware)

	report.HandleFunc("/{user_id}", reportController.SubmitReportHandler).Methods("POST")
}
