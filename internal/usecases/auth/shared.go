package authusecase

import (
	"strings"
	"unicode"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

type TokenService interface {
	GenerateAccessToken(user user.User) (string, error)
	GenerateRefreshToken(user user.User) (string, error)
	ExtractEmailFromRefreshToken(tokenText string) (string, error)
}

type AuthResult struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	ExpiresInSeconds int64
	User             user.User
}

func buildAuthResult(user user.User, tokens TokenService, accessTokenExpirationSeconds int64) (AuthResult, error) {
	accessToken, err := tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	refreshToken, err := tokens.GenerateRefreshToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresInSeconds: accessTokenExpirationSeconds, User: user}, nil
}

func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }

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
