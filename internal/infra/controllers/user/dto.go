package usercontroller

import (
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

type UserResponse struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Role      user.Role `json:"role"`
	CreatedAt time.Time     `json:"createdAt"`
}

type UserSummaryResponse struct {
	ID    int64         `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Role  user.Role `json:"role"`
}

type UpdateUserRoleRequest struct {
	Role user.Role `json:"role"`
}

func ToResponse(user user.User) UserResponse {
	return UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt}
}

func ToSummaryResponse(user *user.User) *UserSummaryResponse {
	if user == nil {
		return nil
	}
	return &UserSummaryResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role}
}
