package userusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type UpdateUserRoleUseCase struct{ users usercontract.UserRepositoryContract }

func NewUpdateUserRoleUseCase(users usercontract.UserRepositoryContract) *UpdateUserRoleUseCase {
	return &UpdateUserRoleUseCase{users: users}
}
func (useCase *UpdateUserRoleUseCase) Execute(ctx context.Context, userID int64, role user.Role) (user.User, error) {
	if !role.IsValidAssignableRole() {
		return user.User{}, domainerrors.Validation("Papel de usuario invalido.")
	}
	updatedUser, found, err := useCase.users.UpdateRole(ctx, userID, role)
	if err != nil {
		return user.User{}, err
	}
	if !found {
		return user.User{}, domainerrors.NotFound("Usuario nao encontrado.")
	}
	return updatedUser, nil
}
