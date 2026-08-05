package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/middleware"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/olabanji12-ojo/church-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportController struct {
	ReportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{ReportService: reportService}
}

// SubmitReportHandler handles a user reporting another user
func (rc *ReportController) SubmitReportHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	reportedIDStr := vars["user_id"]

	var input struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	reporterID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	reportedID, err := primitive.ObjectIDFromHex(reportedIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid reported user ID")
		return
	}

	if err := rc.ReportService.SubmitReport(reporterID, reportedID, input.Reason); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to submit report")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Report submitted successfully"})
}
