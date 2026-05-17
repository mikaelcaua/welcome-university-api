package subjectrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
)

type PostgresSubjectRepositoryImpl struct{ pool *pgxpool.Pool }

func NewPostgresSubjectRepositoryImpl(pool *pgxpool.Pool) *PostgresSubjectRepositoryImpl {
	return &PostgresSubjectRepositoryImpl{pool: pool}
}

func (repository *PostgresSubjectRepositoryImpl) ListByCourse(ctx context.Context, courseID int64) ([]subject.Subject, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM subject WHERE course_id = $1 ORDER BY name ASC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subjects := []subject.Subject{}
	for rows.Next() {
		var subject subject.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, rows.Err()
}

func (repository *PostgresSubjectRepositoryImpl) FindTreeByID(ctx context.Context, subjectID int64) (subject.Subject, bool, error) {
	var subject subject.Subject
	var course course.Course
	var university university.University
	var state state.State
	err := repository.pool.QueryRow(ctx, `
		SELECT
			s.id, s.name,
			c.id, c.name,
			u.id, u.name,
			st.id, st.code, st.name
		FROM subject s
		JOIN course c ON c.id = s.course_id
		JOIN university u ON u.id = c.university_id
		JOIN state st ON st.id = u.state_id
		WHERE s.id = $1
	`, subjectID).Scan(
		&subject.ID, &subject.Name,
		&course.ID, &course.Name,
		&university.ID, &university.Name,
		&state.ID, &state.Code, &state.Name,
	)
	if domainerrors.Is(err, pgx.ErrNoRows) {
		return subject.Subject{}, false, nil
	}
	if err != nil {
		return subject.Subject{}, false, err
	}
	university.State = &state
	course.University = &university
	subject.Course = &course
	return subject, true, nil
}

func (repository *PostgresSubjectRepositoryImpl) NameExistsInCourse(ctx context.Context, name string, courseID int64) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM subject WHERE name = $1 AND course_id = $2)`, name, courseID).Scan(&exists)
	return exists, err
}

func (repository *PostgresSubjectRepositoryImpl) Create(ctx context.Context, subject subject.Subject, courseID int64) (subject.Subject, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO subject (name, course_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, subject.Name, courseID).Scan(&subject.ID, &subject.Name)
	return subject, err
}
