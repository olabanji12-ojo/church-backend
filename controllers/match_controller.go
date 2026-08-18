package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olabanji12-ojo/church-backend/middleware"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/olabanji12-ojo/church-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchController struct {
	SwipeService *services.SwipeService
}

func NewMatchController(swipeService *services.SwipeService) *MatchController {
	return &MatchController{SwipeService: swipeService}
}

// SwipeRightHandler handles liking a user
func (mc *MatchController) SwipeRightHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	targetIDStr := vars["target_user_id"]
	
	actorID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	targetID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	match, err := mc.SwipeService.SwipeRight(actorID, targetID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Swipe failed")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Swipe recorded",
		"data":    match,
	})
}

// SwipeLeftHandler handles passing on a user
func (mc *MatchController) SwipeLeftHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	targetIDStr := vars["target_user_id"]
	
	actorID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	targetID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	match, err := mc.SwipeService.SwipeLeft(actorID, targetID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Swipe failed")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Pass recorded",
		"data":    match,
	})
}

// GetMatchesHandler gets all matched users for the current user
func (mc *MatchController) GetMatchesHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)

	matches, err := mc.SwipeService.GetMatches(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get matches")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Matches retrieved successfully",
		"data":    matches,
	})
}

// GetPendingLikesHandler gets all pending incoming connection requests for the current user
func (mc *MatchController) GetPendingLikesHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(authCtx.UserID)

	likes, err := mc.SwipeService.GetPendingLikes(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get connection requests")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Connection requests retrieved successfully",
		"data":    likes,
	})
}

// UnmatchHandler deletes a match and its messages
func (mc *MatchController) UnmatchHandler(w http.ResponseWriter, r *http.Request) {
	authCtx, err := middleware.GetAuthContextDirect(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	targetIDStr := vars["target_user_id"]
	
	actorID, _ := primitive.ObjectIDFromHex(authCtx.UserID)
	targetID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	err = mc.SwipeService.Unmatch(actorID, targetID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to unmatch")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Unmatched successfully",
	})
}
