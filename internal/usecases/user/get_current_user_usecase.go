package userusecase

import "github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"

type GetCurrentUserUseCase struct{}

func NewGetCurrentUserUseCase() *GetCurrentUserUseCase { return &GetCurrentUserUseCase{} }
func (useCase *GetCurrentUserUseCase) Execute(currentUser user.User) user.User {
	return currentUser
}
