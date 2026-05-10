package services

import (
	"context"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, bool, error)
	FindByID(ctx context.Context, userID int64) (models.User, bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]models.User, error)
	UpdateRole(ctx context.Context, userID int64, role models.Role) (models.User, bool, error)
}

type TokenService interface {
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(user models.User) (string, error)
	ExtractEmailFromRefreshToken(tokenText string) (string, error)
}

type AuthService struct {
	users                        UserRepository
	tokens                       TokenService
	accessTokenExpirationSeconds int64
}

func NewAuthService(users UserRepository, tokens TokenService, accessTokenExpirationSeconds int64) *AuthService {
	return &AuthService{
		users:                       users,
		tokens:                      tokens,
		accessTokenExpirationSeconds: accessTokenExpirationSeconds,
	}
}

func (service *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error) {
	normalizedEmail := normalizeEmail(request.Email)
	if strings.TrimSpace(request.Name) == "" || normalizedEmail == "" {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusBadRequest, "Nome e email sao obrigatorios.")
	}
	if !passwordMeetsComplexityRequirement(request.Password) {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusBadRequest, "A senha deve conter pelo menos 1 letra, 1 numero e 1 caractere especial.")
	}

	exists, err := service.users.EmailExists(ctx, normalizedEmail)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if exists {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusConflict, "Email ja cadastrado.")
	}

	totalUsers, err := service.users.Count(ctx)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	initialRole := models.RoleUser
	if totalUsers == 0 {
		initialRole = models.RoleAdmin
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	user, err := service.users.Create(ctx, models.User{
		Name:         strings.TrimSpace(request.Name),
		Email:        normalizedEmail,
		PasswordHash: string(passwordHash),
		Role:         initialRole,
	})
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return service.buildAuthResponse(user)
}

func (service *AuthService) Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error) {
	user, found, err := service.users.FindByEmail(ctx, normalizeEmail(request.Email))
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if !found || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusUnauthorized, "Credenciais invalidas.")
	}
	return service.buildAuthResponse(user)
}

func (service *AuthService) Refresh(ctx context.Context, request dto.RefreshTokenRequest) (dto.AuthResponse, error) {
	email, err := service.tokens.ExtractEmailFromRefreshToken(request.RefreshToken)
	if err != nil {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusUnauthorized, "Refresh token invalido.")
	}
	user, found, err := service.users.FindByEmail(ctx, email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if !found {
		return dto.AuthResponse{}, httpx.NewHTTPError(http.StatusUnauthorized, "Usuario nao encontrado.")
	}
	return service.buildAuthResponse(user)
}

func (service *AuthService) buildAuthResponse(user models.User) (dto.AuthResponse, error) {
	accessToken, err := service.tokens.GenerateAccessToken(user)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	refreshToken, err := service.tokens.GenerateRefreshToken(user)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresInSeconds: service.accessTokenExpirationSeconds,
		User:             dto.UserToResponse(user),
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func passwordMeetsComplexityRequirement(password string) bool {
	if len(password) < 8 || len(password) > 100 {
		return false
	}
	hasLetter := false
	hasDigit := false
	hasSpecial := false
	for _, character := range password {
		switch {
		case unicode.IsLetter(character):
			hasLetter = true
		case unicode.IsDigit(character):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	return hasLetter && hasDigit && hasSpecial
}
