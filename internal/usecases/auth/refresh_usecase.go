package authusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type RefreshInput struct{ RefreshToken string }

type RefreshUseCase struct {
	users                        usercontract.UserRepositoryContract
	tokens                       TokenService
	accessTokenExpirationSeconds int64
}

func NewRefreshUseCase(users usercontract.UserRepositoryContract, tokens TokenService, accessTokenExpirationSeconds int64) *RefreshUseCase {
	return &RefreshUseCase{users: users, tokens: tokens, accessTokenExpirationSeconds: accessTokenExpirationSeconds}
}
func (useCase *RefreshUseCase) Execute(ctx context.Context, input RefreshInput) (AuthResult, error) {
	email, err := useCase.tokens.ExtractEmailFromRefreshToken(input.RefreshToken)
	if err != nil {
		return AuthResult{}, domainerrors.Unauthorized("Refresh token invalido.")
	}
	user, found, err := useCase.users.FindByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, err
	}
	if !found {
		return AuthResult{}, domainerrors.Unauthorized("Usuario nao encontrado.")
	}
	return buildAuthResult(user, useCase.tokens, useCase.accessTokenExpirationSeconds)
}
