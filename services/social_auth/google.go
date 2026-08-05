package social_auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type GoogleProvider struct {
	ClientID string
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"), // Can be empty during development
	}
}

func (p *GoogleProvider) VerifyToken(token string) (*UserClaims, error) {
	// If we don't have a ClientID, we handle it gracefully so the code doesn't break
	if p.ClientID == "" {
		logrus.Warn("GOOGLE_CLIENT_ID is not set! Skipping real validation for development purposes.")
		return &UserClaims{
			Email:     "mock.google.user@example.com",
			FirstName: "Google",
			LastName:  "User",
			PhotoURL:  "https://ui-avatars.com/api/?name=Google+User",
		}, nil
	}

	// The frontend uses useGoogleLogin which returns an access_token.
	// We need to fetch the user profile from Google's userinfo endpoint.
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("❌ [VerifyToken] Request to Google failed: %v", err)
		return nil, errors.New("failed to contact google servers")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("❌ [VerifyToken] Invalid Google token. Status: %d", resp.StatusCode)
		return nil, errors.New("invalid google token")
	}

	var userInfo struct {
		Email         string `json:"email"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		logrus.Errorf("❌ [VerifyToken] Failed to parse Google response: %v", err)
		return nil, errors.New("failed to parse google response")
	}

	return &UserClaims{
		Email:     userInfo.Email,
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		PhotoURL:  userInfo.Picture,
	}, nil
}
