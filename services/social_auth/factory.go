package social_auth

import "errors"

// GetSocialProvider acts as the Factory for our Social Auth strategies.
// Based on the provider string (e.g., "google", "apple"), it returns the corresponding implementation.
func GetSocialProvider(providerName string) (SocialProvider, error) {
	switch providerName {
	case "google":
		return NewGoogleProvider(), nil
	case "facebook":
		return NewFacebookProvider(), nil
	default:
		return nil, errors.New("unsupported auth provider")
	}
}
