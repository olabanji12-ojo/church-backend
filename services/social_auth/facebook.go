package social_auth

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

type FacebookProvider struct {
	AppID     string
	AppSecret string
}

func NewFacebookProvider() *FacebookProvider {
	return &FacebookProvider{
		AppID:     os.Getenv("FACEBOOK_APP_ID"),
		AppSecret: os.Getenv("FACEBOOK_APP_SECRET"),
	}
}

func (p *FacebookProvider) VerifyToken(token string) (*UserClaims, error) {
	// Graceful fallback for missing keys
	if p.AppID == "" {
		logrus.Warn("FACEBOOK_APP_ID is not set! Skipping real validation for development purposes.")
		return &UserClaims{
			Email:     "mock.facebook.user@example.com",
			FirstName: "Facebook",
			LastName:  "User",
			PhotoURL:  "https://ui-avatars.com/api/?name=Facebook+User",
		}, nil
	}

	// TODO: Make HTTP request to graph.facebook.com/me?access_token=token
	return nil, errors.New("real facebook verification logic not fully implemented yet")
}
