package usercontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

type UserRepositoryContract interface {
	Create(ctx context.Context, user user.User) (user.User, error)
	FindByEmail(ctx context.Context, email string) (user.User, bool, error)
	FindByID(ctx context.Context, userID int64) (user.User, bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]user.User, error)
	UpdateRole(ctx context.Context, userID int64, role user.Role) (user.User, bool, error)
}
