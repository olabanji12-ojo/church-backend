package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/services"
	"github.com/olabanji12-ojo/church-backend/utils"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// REGISTER HANDLER
func (ac *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logrus.Errorf("❌ [RegisterHandler] Failed to parse JSON: %v", err)
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation (can be expanded later)
	if input.Email == "" || input.PasswordHash == "" {
		utils.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	newUser, err := ac.AuthService.RegisterUser(input)
	if err != nil {
		logrus.Errorf("❌ [RegisterHandler] Registration failed: %v", err)
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate JWT token for the new user
	token, err := utils.GenerateJWT(newUser.ID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: func() http.SameSite {
			if os.Getenv("ENVIRONMENT") == "production" {
				return http.SameSiteNoneMode
			}
			return http.SameSiteStrictMode
		}(),
	})

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successful",
		"data": map[string]interface{}{
			"user":  newUser,
			"token": token,
		},
	})
}

// LOGIN HANDLER
func (ac *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, user, err := ac.AuthService.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: func() http.SameSite {
			if os.Getenv("ENVIRONMENT") == "production" {
				return http.SameSiteNoneMode
			}
			return http.SameSiteStrictMode
		}(),
	})

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}

// LOGOUT HANDLER
func (ac *AuthController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: func() http.SameSite {
			if os.Getenv("ENVIRONMENT") == "production" {
				return http.SameSiteNoneMode
			}
			return http.SameSiteStrictMode
		}(),
	})
	utils.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// SOCIAL AUTH HANDLER
func (ac *AuthController) SocialAuthHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Provider string `json:"provider"`
		Token    string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Provider == "" || input.Token == "" {
		utils.Error(w, http.StatusBadRequest, "Provider and token are required")
		return
	}

	token, user, err := ac.AuthService.HandleSocialLogin(input.Provider, input.Token)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: func() http.SameSite {
			if os.Getenv("ENVIRONMENT") == "production" {
				return http.SameSiteNoneMode
			}
			return http.SameSiteStrictMode
		}(),
	})

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Social login successful",
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}

// GUEST LOGIN HANDLER
func (ac *AuthController) GuestLoginHandler(w http.ResponseWriter, r *http.Request) {
	token, user, err := ac.AuthService.GuestLogin()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create guest account")
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: func() http.SameSite {
			if os.Getenv("ENVIRONMENT") == "production" {
				return http.SameSiteNoneMode
			}
			return http.SameSiteStrictMode
		}(),
	})

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Guest login successful",
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}
