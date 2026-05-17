package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikaelcaua/welcome-university-api/internal/infra/config"
)

func Connect(ctx context.Context, appConfig config.Config) (*pgxpool.Pool, error) {
	dsn := appConfig.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@db:5432/%s?sslmode=disable",
			url.QueryEscape(appConfig.DatabaseUser),
			url.QueryEscape(appConfig.DatabasePassword),
			url.QueryEscape(appConfig.DatabaseName),
		)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
