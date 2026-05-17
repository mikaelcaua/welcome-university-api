package universityrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
)

type PostgresUniversityRepositoryImpl struct{ pool *pgxpool.Pool }

func NewPostgresUniversityRepositoryImpl(pool *pgxpool.Pool) *PostgresUniversityRepositoryImpl {
	return &PostgresUniversityRepositoryImpl{pool: pool}
}

func (repository *PostgresUniversityRepositoryImpl) ListByState(ctx context.Context, stateID int64) ([]university.University, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM university WHERE state_id = $1 ORDER BY name ASC`, stateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	universities := []university.University{}
	for rows.Next() {
		var university university.University
		if err := rows.Scan(&university.ID, &university.Name); err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	return universities, rows.Err()
}

func (repository *PostgresUniversityRepositoryImpl) FindByID(ctx context.Context, universityID int64) (university.University, bool, error) {
	var university university.University
	err := repository.pool.QueryRow(ctx, `SELECT id, name FROM university WHERE id = $1`, universityID).Scan(&university.ID, &university.Name)
	if domainerrors.Is(err, pgx.ErrNoRows) {
		return university.University{}, false, nil
	}
	return university, err == nil, err
}

func (repository *PostgresUniversityRepositoryImpl) NameExistsInState(ctx context.Context, name string, stateID int64) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM university WHERE name = $1 AND state_id = $2)`, name, stateID).Scan(&exists)
	return exists, err
}

func (repository *PostgresUniversityRepositoryImpl) Create(ctx context.Context, university university.University, stateID int64) (university.University, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO university (name, state_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, university.Name, stateID).Scan(&university.ID, &university.Name)
	return university, err
}
