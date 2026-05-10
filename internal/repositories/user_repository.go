package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (repository *PostgresUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO app_user (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, created_at, updated_at
	`, user.Name, strings.ToLower(user.Email), user.PasswordHash, user.Role).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (repository *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (models.User, bool, error) {
	var user models.User
	err := repository.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		WHERE lower(email) = lower($1)
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, false, nil
	}
	return user, err == nil, err
}

func (repository *PostgresUserRepository) FindByID(ctx context.Context, userID int64) (models.User, bool, error) {
	var user models.User
	err := repository.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, false, nil
	}
	return user, err == nil, err
}

func (repository *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM app_user WHERE lower(email) = lower($1))`, email).Scan(&exists)
	return exists, err
}

func (repository *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := repository.pool.QueryRow(ctx, `SELECT count(*) FROM app_user`).Scan(&total)
	return total, err
}

func (repository *PostgresUserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM app_user
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repository *PostgresUserRepository) UpdateRole(ctx context.Context, userID int64, role models.Role) (models.User, bool, error) {
	var user models.User
	err := repository.pool.QueryRow(ctx, `
		UPDATE app_user
		SET role = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, name, email, password_hash, role, created_at, updated_at
	`, userID, role).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, false, nil
	}
	return user, err == nil, err
}
