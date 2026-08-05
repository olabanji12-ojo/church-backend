package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/olabanji12-ojo/church-backend/utils"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

// AuthContext stores the data we'll inject into the request context
type AuthContext struct {
	UserID string
}

// typed context key to prevent collisions
type contextKey string

const authKey contextKey = "auth"

// GetAuthContext is a helper function that controllers can use for safe context retrieval
func GetAuthContext(r *http.Request) (*AuthContext, bool) {
	authValue := r.Context().Value(authKey)
	if authValue == nil {
		logrus.Warn("Auth context not found in request")
		return nil, false
	}

	authCtx, ok := authValue.(AuthContext)
	if !ok {
		logrus.Warn("Auth context type assertion failed")
		return nil, false
	}

	return &authCtx, true
}

// GetAuthContextDirect is an alternative helper matching the Car Wash pattern
func GetAuthContextDirect(r *http.Request) (AuthContext, error) {
	authValue := r.Context().Value(authKey)
	if authValue == nil {
		return AuthContext{}, fmt.Errorf("authentication context not found")
	}

	authCtx, ok := authValue.(AuthContext)
	if !ok {
		return AuthContext{}, fmt.Errorf("invalid authentication context type")
	}

	return authCtx, nil
}

// AuthMiddleware checks for token, validates it, and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Try Authorization header first (Bearer Token)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Try HttpOnly cookie
			cookie, err := r.Cookie("jwt")
			if err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				// Try query param (Useful for WebSockets later)
				tokenString = r.URL.Query().Get("auth_token")
			}
		}

		if tokenString == "" {
			logrus.Warn("No auth token provided")
			http.Error(w, `{"error": "Missing auth token"}`, http.StatusUnauthorized)
			return
		}

		// Validate the token
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil || userID == "" {
			logrus.WithError(err).Warn("Token invalid or expired")
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		authCtx := AuthContext{
			UserID: userID,
		}

		// Inject AuthContext into request context
		ctx := context.WithValue(r.Context(), authKey, authCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Cors sets up CORS rules strictly following the Car Wash app pattern
func Cors() *cors.Cors {
	origins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:5173", // Vite default
		"http://127.0.0.1:5173",
		"https://covenant-frontend-liard.vercel.app", // Production Vercel URL
	}

	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	if envOrigins != "" {
		for _, o := range strings.Split(envOrigins, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}

	return cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		Debug:            os.Getenv("ENVIRONMENT") != "production",
	})
}

// Secure sets up basic security headers strictly following the Car Wash app pattern
func Secure() *secure.Secure {
	options := secure.Options{
		BrowserXssFilter:     true,
		ContentTypeNosniff:   true,
		FrameDeny:            false,
		SSLForceHost:         false,
		STSIncludeSubdomains: true,
		STSPreload:           true,
	}
	if os.Getenv("ENVIRONMENT") != "production" {
		options.IsDevelopment = true
	}
	return secure.New(options)
}
