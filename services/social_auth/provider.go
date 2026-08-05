package social_auth

// UserClaims represents the standardized data we extract from any social provider's token
type UserClaims struct {
	Email     string
	FirstName string
	LastName  string
	PhotoURL  string
}

// SocialProvider is the unified interface all social login strategies must implement
type SocialProvider interface {
	VerifyToken(token string) (*UserClaims, error)
}
