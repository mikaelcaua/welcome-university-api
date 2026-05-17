package staterepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
)

type PostgresStateRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewPostgresStateRepositoryImpl(pool *pgxpool.Pool) *PostgresStateRepositoryImpl {
	return &PostgresStateRepositoryImpl{pool: pool}
}

func (repository *PostgresStateRepositoryImpl) List(ctx context.Context) ([]state.State, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, code, name FROM state ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := []state.State{}
	for rows.Next() {
		var state state.State
		if err := rows.Scan(&state.ID, &state.Code, &state.Name); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (repository *PostgresStateRepositoryImpl) FindByID(ctx context.Context, stateID int64) (state.State, bool, error) {
	return repository.find(ctx, `SELECT id, code, name FROM state WHERE id = $1`, stateID)
}

func (repository *PostgresStateRepositoryImpl) FindByCode(ctx context.Context, code string) (state.State, bool, error) {
	return repository.find(ctx, `SELECT id, code, name FROM state WHERE code = $1`, code)
}

func (repository *PostgresStateRepositoryImpl) Create(ctx context.Context, state state.State) (state.State, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO state (code, name)
		VALUES ($1, $2)
		RETURNING id, code, name
	`, state.Code, state.Name).Scan(&state.ID, &state.Code, &state.Name)
	return state, err
}

func (repository *PostgresStateRepositoryImpl) find(ctx context.Context, query string, arg any) (state.State, bool, error) {
	var state state.State
	err := repository.pool.QueryRow(ctx, query, arg).Scan(&state.ID, &state.Code, &state.Name)
	if domainerrors.Is(err, pgx.ErrNoRows) {
		return state.State{}, false, nil
	}
	return state, err == nil, err
}
