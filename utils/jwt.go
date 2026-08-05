package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/olabanji12-ojo/church-backend/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateJWT creates a new JSON Web Token for the given User ID
func GenerateJWT(userID primitive.ObjectID) (string, error) {
	secretKey := config.GetEnv("JWT_SECRET", "super_secret_jwt_key_for_dev")

	// Token expires in 7 days
	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateJWT parses the token and returns the user ID if valid
func ValidateJWT(tokenString string) (string, error) {
	secretKey := config.GetEnv("JWT_SECRET", "super_secret_jwt_key_for_dev")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, exists := claims["user_id"].(string)
		if !exists {
			return "", errors.New("user_id not found in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
