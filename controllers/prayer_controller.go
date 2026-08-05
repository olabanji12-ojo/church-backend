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

type PrayerController struct {
	PrayerService *services.PrayerService
}

func NewPrayerController(prayerService *services.PrayerService) *PrayerController {
	return &PrayerController{PrayerService: prayerService}
}

// GetPrayersHandler fetches the community prayer feed
func (pc *PrayerController) GetPrayersHandler(w http.ResponseWriter, r *http.Request) {
	// Authentication is optional for viewing if we want, but let's assume protected
	_, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	prayers, err := pc.PrayerService.GetPrayerFeed()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load prayers")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": prayers,
	})
}

// PostPrayerHandler creates a new prayer request
func (pc *PrayerController) PostPrayerHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	authorID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	if err := pc.PrayerService.PostPrayer(authorID, input.Content); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to post prayer")
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]string{"message": "Prayer posted successfully"})
}

// AddAmenHandler allows a user to "Amen" a prayer
func (pc *PrayerController) AddAmenHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	prayerIDStr := mux.Vars(r)["id"]
	prayerID, err := primitive.ObjectIDFromHex(prayerIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid prayer ID")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)

	// We need to add this method to PrayerService
	if err := pc.PrayerService.SayAmen(prayerID, userID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to add Amen")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Amen added successfully"})
}
