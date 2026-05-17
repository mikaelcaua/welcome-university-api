package authcontroller

import "github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken      string                      `json:"accessToken"`
	RefreshToken     string                      `json:"refreshToken"`
	TokenType        string                      `json:"tokenType"`
	ExpiresInSeconds int64                       `json:"expiresInSeconds"`
	User             usercontroller.UserResponse `json:"user"`
}
