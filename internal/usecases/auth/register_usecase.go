package authusecase

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUseCase struct {
	users                        usercontract.UserRepositoryContract
	tokens                       TokenService
	accessTokenExpirationSeconds int64
}

func NewRegisterUseCase(users usercontract.UserRepositoryContract, tokens TokenService, accessTokenExpirationSeconds int64) *RegisterUseCase {
	return &RegisterUseCase{users: users, tokens: tokens, accessTokenExpirationSeconds: accessTokenExpirationSeconds}
}
func (useCase *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (AuthResult, error) {
	user := user.User{Name: strings.TrimSpace(input.Name), Email: normalizeEmail(input.Email)}
	if !user.HasValidIdentity() {
		return AuthResult{}, domainerrors.Validation("Nome e email sao obrigatorios.")
	}
	if !passwordMeetsComplexityRequirement(input.Password) {
		return AuthResult{}, domainerrors.Validation("A senha deve conter pelo menos 1 letra, 1 numero e 1 caractere especial.")
	}
	exists, err := useCase.users.EmailExists(ctx, user.Email)
	if err != nil {
		return AuthResult{}, err
	}
	if exists {
		return AuthResult{}, domainerrors.Conflict("Email ja cadastrado.")
	}
	totalUsers, err := useCase.users.Count(ctx)
	if err != nil {
		return AuthResult{}, err
	}
	user.Role = user.RoleUser
	if totalUsers == 0 {
		user.Role = user.RoleAdmin
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}
	user.PasswordHash = string(passwordHash)
	user, err = useCase.users.Create(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}
	return buildAuthResult(user, useCase.tokens, useCase.accessTokenExpirationSeconds)
}
