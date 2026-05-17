package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type JWTService struct {
	secret                         []byte
	accessTokenExpirationDuration  time.Duration
	refreshTokenExpirationDuration time.Duration
}

func NewJWTService(secret string, accessTokenExpirationSeconds int64, refreshTokenExpirationSeconds int64) *JWTService {
	return &JWTService{
		secret:                         []byte(secret),
		accessTokenExpirationDuration:  time.Duration(accessTokenExpirationSeconds) * time.Second,
		refreshTokenExpirationDuration: time.Duration(refreshTokenExpirationSeconds) * time.Second,
	}
}

func (service *JWTService) GenerateAccessToken(user user.User) (string, error) {
	return service.generateToken(user, AccessTokenType, service.accessTokenExpirationDuration)
}

func (service *JWTService) GenerateRefreshToken(user user.User) (string, error) {
	return service.generateToken(user, RefreshTokenType, service.refreshTokenExpirationDuration)
}

func (service *JWTService) ExtractEmailFromAccessToken(tokenText string) (string, error) {
	return service.extractEmail(tokenText, AccessTokenType)
}

func (service *JWTService) ExtractEmailFromRefreshToken(tokenText string) (string, error) {
	return service.extractEmail(tokenText, RefreshTokenType)
}

func (service *JWTService) generateToken(user user.User, tokenType string, expirationDuration time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.Email,
		"uid":  user.ID,
		"role": string(user.Role),
		"type": tokenType,
		"iat":  now.Unix(),
		"exp":  now.Add(expirationDuration).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(service.secret)
}

func (service *JWTService) extractEmail(tokenText string, expectedTokenType string) (string, error) {
	token, err := jwt.Parse(tokenText, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, domainerrors.New("invalid signing method")
		}
		return service.secret, nil
	})
	if err != nil || !token.Valid {
		return "", domainerrors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", domainerrors.New("invalid claims")
	}
	tokenType, _ := claims["type"].(string)
	if tokenType != expectedTokenType {
		return "", domainerrors.New("invalid token type")
	}
	email, _ := claims["sub"].(string)
	if email == "" {
		return "", domainerrors.New("missing email")
	}
	return email, nil
}
