package controllers

import (
	"encoding/json"
	"net/http"
	"log"
	"time"

	"github.com/olabanji12-ojo/church-backend/middleware"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/olabanji12-ojo/church-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	ProfileService *services.ProfileService
	SwipeService   *services.SwipeService
}

func NewUserController(profileService *services.ProfileService, swipeService *services.SwipeService) *UserController {
	return &UserController{
		ProfileService: profileService,
		SwipeService:   swipeService,
	}
}

// GetDiscoveryFeedHandler fetches potential matches
func (uc *UserController) GetDiscoveryFeedHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)

	feed, err := uc.SwipeService.GetDiscoveryFeed(userID)
	if err != nil {
		log.Printf("ERROR in GetDiscoveryFeed: %v", err)
		utils.Error(w, http.StatusInternalServerError, "Failed to load discovery feed")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": feed,
	})
}

// UpdateFaithProfileHandler
func (uc *UserController) UpdateFaithProfileHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Denomination string `json:"denomination"`
		ChurchFreq   string `json:"church_freq"`
		PrayerFreq   string `json:"prayer_freq"`
		BibleFreq    string `json:"bible_freq"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	if err := uc.ProfileService.UpdateFaithProfile(userID, input.Denomination, input.ChurchFreq, input.PrayerFreq, input.BibleFreq); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Faith profile updated"})
}

// GetCurrentUserHandler fetches the current logged-in user profile
func (uc *UserController) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	user, err := uc.ProfileService.GetUserProfile(userID)
	if err != nil {
		log.Printf("ERROR in GetCurrentUserHandler: %v", err)
		utils.Error(w, http.StatusNotFound, "User not found")
		return
	}

	user.PasswordHash = "" // Security: Omit password hash
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// UpdateCurrentUserHandler updates the current logged-in user profile
func (uc *UserController) UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// We clean or adjust keys as needed for MongoDB.
	delete(updateData, "id")
	delete(updateData, "email")
	delete(updateData, "password_hash")
	delete(updateData, "created_at")

	// Convert keys to bson.M
	updateBson := make(primitive.M)
	for k, v := range updateData {
		if k == "dob" {
			if dobStr, ok := v.(string); ok && dobStr != "" {
				if parsedTime, err := time.Parse(time.RFC3339, dobStr); err == nil {
					updateBson[k] = parsedTime
				} else if parsedTime, err := time.Parse("2006-01-02", dobStr); err == nil {
					updateBson[k] = parsedTime
				} else {
					log.Printf("⚠️ Failed to parse dob string: %s", dobStr)
					updateBson[k] = v
				}
			} else {
				updateBson[k] = v
			}
		} else {
			updateBson[k] = v
		}
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	if err := uc.ProfileService.UpdateUserProfile(userID, updateBson); err != nil {
		log.Printf("ERROR in UpdateCurrentUserHandler: %v", err)
		utils.Error(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Fetch updated user to return it
	user, err := uc.ProfileService.GetUserProfile(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve updated profile")
		return
	}

	user.PasswordHash = ""
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// SavePushToken saves the Firebase FCM token for the currently logged in user
func (uc *UserController) SavePushToken(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)

	var body struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if body.Token == "" {
		utils.Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	// Update user profile using ProfileService
	err = uc.ProfileService.UpdateUserProfile(userID, bson.M{"push_token": body.Token})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save push token")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Push token saved successfully"})
}

// GetScenarioQuestionsHandler returns all predefined Covenant scenario questions
func (uc *UserController) GetScenarioQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	scenarioSvc := services.NewScenarioService()
	questions := scenarioSvc.GetPredefinedScenarios()
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"data": questions,
	})
}

// AnswerScenarioHandler saves a user's answer to a scenario question
func (uc *UserController) AnswerScenarioHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		QuestionID string `json:"question_id"`
		OptionID   string `json:"option_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.QuestionID == "" || req.OptionID == "" {
		utils.Error(w, http.StatusBadRequest, "Invalid question or option ID")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	updateKey := "scenario_answers." + req.QuestionID

	if err := uc.ProfileService.UpdateUserProfile(userID, bson.M{updateKey: req.OptionID}); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save scenario answer")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Answer recorded successfully"})
}

