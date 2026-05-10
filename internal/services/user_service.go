package services

import (
	"context"
	"net/http"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type UserService struct {
	users UserRepository
}

func NewUserService(users UserRepository) *UserService {
	return &UserService{users: users}
}

func (service *UserService) Me(currentUser models.User) dto.UserResponse {
	return dto.UserToResponse(currentUser)
}

func (service *UserService) ListAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := service.users.List(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, dto.UserToResponse(user))
	}
	return response, nil
}

func (service *UserService) UpdateRole(ctx context.Context, userID int64, role models.Role) (dto.UserResponse, error) {
	if !role.IsValidAssignableRole() {
		return dto.UserResponse{}, httpx.NewHTTPError(http.StatusBadRequest, "Papel de usuario invalido.")
	}
	updatedUser, found, err := service.users.UpdateRole(ctx, userID, role)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !found {
		return dto.UserResponse{}, httpx.NewHTTPError(http.StatusNotFound, "Usuario nao encontrado.")
	}
	return dto.UserToResponse(updatedUser), nil
}
