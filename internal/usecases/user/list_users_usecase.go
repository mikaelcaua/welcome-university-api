package userusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/user"
)

type ListUsersUseCase struct{ users usercontract.UserRepositoryContract }

func NewListUsersUseCase(users usercontract.UserRepositoryContract) *ListUsersUseCase {
	return &ListUsersUseCase{users: users}
}
func (useCase *ListUsersUseCase) Execute(ctx context.Context) ([]user.User, error) {
	return useCase.users.List(ctx)
}
