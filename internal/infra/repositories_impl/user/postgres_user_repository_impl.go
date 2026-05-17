package userrepository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

type PostgresUserRepositoryImpl struct{ pool *pgxpool.Pool }

func NewPostgresUserRepositoryImpl(pool *pgxpool.Pool) *PostgresUserRepositoryImpl {
	return &PostgresUserRepositoryImpl{pool: pool}
}

func (repository *PostgresUserRepositoryImpl) Create(ctx context.Context, user user.User) (user.User, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO app_user (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, created_at, updated_at
	`, user.Name, strings.ToLower(user.Email), user.PasswordHash, user.Role).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (repository *PostgresUserRepositoryImpl) FindByEmail(ctx context.Context, email string) (user.User, bool, error) {
	return repository.find(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		WHERE lower(email) = lower($1)
	`, email)
}

func (repository *PostgresUserRepositoryImpl) FindByID(ctx context.Context, userID int64) (user.User, bool, error) {
	return repository.find(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		WHERE id = $1
	`, userID)
}

func (repository *PostgresUserRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM app_user WHERE lower(email) = lower($1))`, email).Scan(&exists)
	return exists, err
}

func (repository *PostgresUserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var total int64
	err := repository.pool.QueryRow(ctx, `SELECT count(*) FROM app_user`).Scan(&total)
	return total, err
}

func (repository *PostgresUserRepositoryImpl) List(ctx context.Context) ([]user.User, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []user.User{}
	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repository *PostgresUserRepositoryImpl) UpdateRole(ctx context.Context, userID int64, role user.Role) (user.User, bool, error) {
	var foundUser user.User
	err := repository.pool.QueryRow(ctx, `
		UPDATE app_user
		SET role = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, name, email, password_hash, role, created_at, updated_at
	`, userID, role).Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.PasswordHash, &foundUser.Role, &foundUser.CreatedAt, &foundUser.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, false, nil
	}
	return foundUser, err == nil, err
}

func (repository *PostgresUserRepositoryImpl) find(ctx context.Context, query string, arg any) (user.User, bool, error) {
	var foundUser user.User
	err := repository.pool.QueryRow(ctx, query, arg).Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.PasswordHash, &foundUser.Role, &foundUser.CreatedAt, &foundUser.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, false, nil
	}
	return foundUser, err == nil, err
}
