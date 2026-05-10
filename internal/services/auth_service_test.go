package services

import (
	"context"
	"testing"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

func TestRegisterFirstUserBecomesAdmin(t *testing.T) {
	userRepository := newFakeUserRepository()
	tokenService := fakeTokenService{}
	service := NewAuthService(userRepository, tokenService, 900)

	response, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "Mikael",
		Email:    "MIKAEL@example.com",
		Password: "Senha@123",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if response.User.Role != models.RoleAdmin {
		t.Fatalf("expected first user to be ADMIN, got %s", response.User.Role)
	}
	if response.User.Email != "mikael@example.com" {
		t.Fatalf("expected normalized email, got %s", response.User.Email)
	}
	if response.AccessToken == "" || response.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}
}

func TestRegisterRejectsWeakPassword(t *testing.T) {
	service := NewAuthService(newFakeUserRepository(), fakeTokenService{}, 900)

	_, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "Mikael",
		Email:    "mikael@example.com",
		Password: "senhafraca",
	})
	if err == nil {
		t.Fatal("expected weak password error")
	}
}

type fakeTokenService struct{}

func (fakeTokenService) GenerateAccessToken(user models.User) (string, error) {
	return "access:" + user.Email, nil
}

func (fakeTokenService) GenerateRefreshToken(user models.User) (string, error) {
	return "refresh:" + user.Email, nil
}

func (fakeTokenService) ExtractEmailFromRefreshToken(tokenText string) (string, error) {
	if len(tokenText) <= len("refresh:") {
		return "", nil
	}
	return tokenText[len("refresh:"):], nil
}

type fakeUserRepository struct {
	usersByID    map[int64]models.User
	usersByEmail map[string]models.User
	nextID       int64
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByID:    map[int64]models.User{},
		usersByEmail: map[string]models.User{},
		nextID:       1,
	}
}

func (repository *fakeUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	user.ID = repository.nextID
	repository.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	repository.usersByID[user.ID] = user
	repository.usersByEmail[user.Email] = user
	return user, nil
}

func (repository *fakeUserRepository) FindByEmail(ctx context.Context, email string) (models.User, bool, error) {
	user, found := repository.usersByEmail[email]
	return user, found, nil
}

func (repository *fakeUserRepository) FindByID(ctx context.Context, userID int64) (models.User, bool, error) {
	user, found := repository.usersByID[userID]
	return user, found, nil
}

func (repository *fakeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	_, found := repository.usersByEmail[email]
	return found, nil
}

func (repository *fakeUserRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(repository.usersByID)), nil
}

func (repository *fakeUserRepository) List(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0, len(repository.usersByID))
	for _, user := range repository.usersByID {
		users = append(users, user)
	}
	return users, nil
}

func (repository *fakeUserRepository) UpdateRole(ctx context.Context, userID int64, role models.Role) (models.User, bool, error) {
	user, found := repository.usersByID[userID]
	if !found {
		return models.User{}, false, nil
	}
	user.Role = role
	repository.usersByID[userID] = user
	repository.usersByEmail[user.Email] = user
	return user, true, nil
}
