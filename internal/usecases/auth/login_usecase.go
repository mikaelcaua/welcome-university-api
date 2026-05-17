package authusecase

import (
	"context"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	users                        usercontract.UserRepositoryContract
	tokens                       TokenService
	accessTokenExpirationSeconds int64
}

func NewLoginUseCase(users usercontract.UserRepositoryContract, tokens TokenService, accessTokenExpirationSeconds int64) *LoginUseCase {
	return &LoginUseCase{users: users, tokens: tokens, accessTokenExpirationSeconds: accessTokenExpirationSeconds}
}
func (useCase *LoginUseCase) Execute(ctx context.Context, input LoginInput) (AuthResult, error) {
	user, found, err := useCase.users.FindByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		return AuthResult{}, err
	}
	if !found || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return AuthResult{}, domainerrors.Unauthorized("Credenciais invalidas.")
	}
	return buildAuthResult(user, useCase.tokens, useCase.accessTokenExpirationSeconds)
}
