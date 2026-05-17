package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	commands := []string{
		`CREATE TABLE IF NOT EXISTS app_user (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE IF NOT EXISTS state (
			id BIGSERIAL PRIMARY KEY,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS university (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			state_id BIGINT REFERENCES state(id)
		)`,
		`CREATE TABLE IF NOT EXISTS course (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			university_id BIGINT REFERENCES university(id)
		)`,
		`CREATE TABLE IF NOT EXISTS subject (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			course_id BIGINT REFERENCES course(id)
		)`,
		`CREATE TABLE IF NOT EXISTS exam (
			id BIGSERIAL PRIMARY KEY,
			name TEXT,
			exam_year INTEGER NOT NULL,
			semester INTEGER NOT NULL,
			period_label TEXT,
			pdf_url TEXT NOT NULL,
			storage_key TEXT,
			file_hash TEXT UNIQUE,
			type TEXT,
			status TEXT,
			subject_id BIGINT REFERENCES subject(id),
			uploaded_by_id BIGINT REFERENCES app_user(id),
			reviewed_by_id BIGINT REFERENCES app_user(id),
			review_note TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			reviewed_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ DEFAULT now()
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_university_name_state ON university (name, state_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_course_name_university ON course (name, university_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_subject_name_course ON subject (name, course_id)`,
		`CREATE INDEX IF NOT EXISTS idx_exam_status_subject_created ON exam (status, subject_id, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_exam_uploaded_status ON exam (uploaded_by_id, status)`,
	}

	for _, command := range commands {
		if _, err := pool.Exec(ctx, command); err != nil {
			return err
		}
	}
	return nil
}
