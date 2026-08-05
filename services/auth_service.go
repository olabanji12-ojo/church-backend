package services

import (
	"errors"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/olabanji12-ojo/church-backend/repositories"
	"github.com/olabanji12-ojo/church-backend/services/social_auth"
	"github.com/olabanji12-ojo/church-backend/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	userRepository   *repositories.UserRepository
	embeddingService *EmbeddingService
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository:   userRepository,
		embeddingService: NewEmbeddingService(),
	}
}

// TriggerEmbeddingUpdate generates the user embedding asynchronously in the background.
func (as *AuthService) TriggerEmbeddingUpdate(userID primitive.ObjectID) {
	go func() {
		user, err := as.userRepository.FindUserByID(userID)
		if err != nil {
			return
		}

		if user.IsGuest {
			return
		}

		text := as.embeddingService.GenerateUserText(user)
		embedding, err := as.embeddingService.GetEmbedding(text)
		if err != nil {
			return
		}

		_ = as.userRepository.UpdateUserByID(userID, bson.M{"profile_embedding": embedding})
	}()
}

// RegisterUser hashes the password and saves the user
func (as *AuthService) RegisterUser(input models.User) (*models.User, error) {
	// 1. Check if user with email already exists
	existing, _ := as.userRepository.FindUserByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2. Hash the password
	hashedPassword, err := utils.HashPassword(input.PasswordHash) // The plain password comes in via PasswordHash field temporarily
	if err != nil {
		logrus.Error("Error hashing password: ", err)
		return nil, err
	}

	input.PasswordHash = hashedPassword
	input.IsVerified = false

	// 3. Insert into DB via Repository
	err = as.userRepository.CreateUser(&input)
	if err != nil {
		return nil, err
	}
	as.TriggerEmbeddingUpdate(input.ID)

	input.PasswordHash = "" // Clear before returning
	return &input, nil
}

// LoginUser verifies credentials and returns a JWT
func (as *AuthService) LoginUser(email, password string) (string, *models.User, error) {
	// 1. Find user by email
	user, err := as.userRepository.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	// 2. Check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid password")
	}

	// 3. Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		logrus.Error("Error generating token: ", err)
		return "", nil, err
	}

	user.PasswordHash = "" // Clear before returning
	return token, user, nil
}

// VerifyEmail (Placeholder for OTP logic)
func (as *AuthService) VerifyEmail(email, token string) error {
	// TODO: Implement OTP checking using Redis or User collection
	return nil
}

// HandleSocialLogin uses the Auth Factory to verify the token, then finds or creates the user.
func (as *AuthService) HandleSocialLogin(providerName, token string) (string, *models.User, error) {
	// 1. Get the provider from Factory
	provider, err := social_auth.GetSocialProvider(providerName)
	if err != nil {
		return "", nil, err
	}

	// 2. Verify the token
	claims, err := provider.VerifyToken(token)
	if err != nil {
		return "", nil, err
	}

	// 3. Find user by email
	user, _ := as.userRepository.FindUserByEmail(claims.Email)

	// 4. If user doesn't exist, create them
	if user == nil {
		newUser := models.User{
			Email:      claims.Email,
			FirstName:  claims.FirstName,
			LastName:   claims.LastName,
			IsVerified: true, // Social accounts are implicitly verified
		}
		if len(claims.PhotoURL) > 0 {
			newUser.Photos = []string{claims.PhotoURL}
		}

		err = as.userRepository.CreateUser(&newUser)
		if err != nil {
			return "", nil, err
		}
		user = &newUser
		as.TriggerEmbeddingUpdate(user.ID)
	}

	// 5. Generate JWT token
	jwtToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	return jwtToken, user, nil
}

// GuestLogin creates a temporary shadow user and returns their JWT.
func (as *AuthService) GuestLogin() (string, *models.User, error) {
	// Create a shadow user
	shadowUser := models.User{
		IsGuest:   true,
		FirstName: "Guest",
		LastName:  "User",
		Email:     "guest_" + utils.GenerateRandomString(10) + "@churchmatch.temp", // Generate random email
	}

	err := as.userRepository.CreateUser(&shadowUser)
	if err != nil {
		return "", nil, err
	}

	jwtToken, err := utils.GenerateJWT(shadowUser.ID)
	if err != nil {
		return "", nil, err
	}

	return jwtToken, &shadowUser, nil
}

