package dto

import (
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

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
	AccessToken      string       `json:"accessToken"`
	RefreshToken     string       `json:"refreshToken"`
	TokenType        string       `json:"tokenType"`
	ExpiresInSeconds int64        `json:"expiresInSeconds"`
	User             UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UserSummaryResponse struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

type UpdateUserRoleRequest struct {
	Role models.Role `json:"role"`
}

func UserToResponse(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func UserToSummaryResponse(user *models.User) *UserSummaryResponse {
	if user == nil {
		return nil
	}
	return &UserSummaryResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
