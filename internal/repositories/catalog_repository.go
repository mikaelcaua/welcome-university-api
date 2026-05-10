package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type PostgresCatalogRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCatalogRepository(pool *pgxpool.Pool) *PostgresCatalogRepository {
	return &PostgresCatalogRepository{pool: pool}
}

func (repository *PostgresCatalogRepository) ListStates(ctx context.Context) ([]models.State, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, code, name FROM state ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := []models.State{}
	for rows.Next() {
		var state models.State
		if err := rows.Scan(&state.ID, &state.Code, &state.Name); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (repository *PostgresCatalogRepository) FindStateByID(ctx context.Context, stateID int64) (models.State, bool, error) {
	var state models.State
	err := repository.pool.QueryRow(ctx, `SELECT id, code, name FROM state WHERE id = $1`, stateID).Scan(&state.ID, &state.Code, &state.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.State{}, false, nil
	}
	return state, err == nil, err
}

func (repository *PostgresCatalogRepository) FindStateByCode(ctx context.Context, code string) (models.State, bool, error) {
	var state models.State
	err := repository.pool.QueryRow(ctx, `SELECT id, code, name FROM state WHERE code = $1`, code).Scan(&state.ID, &state.Code, &state.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.State{}, false, nil
	}
	return state, err == nil, err
}

func (repository *PostgresCatalogRepository) CreateState(ctx context.Context, state models.State) (models.State, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO state (code, name)
		VALUES ($1, $2)
		RETURNING id, code, name
	`, state.Code, state.Name).Scan(&state.ID, &state.Code, &state.Name)
	return state, err
}

func (repository *PostgresCatalogRepository) ListUniversitiesByState(ctx context.Context, stateID int64) ([]models.University, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM university WHERE state_id = $1 ORDER BY name ASC`, stateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	universities := []models.University{}
	for rows.Next() {
		var university models.University
		if err := rows.Scan(&university.ID, &university.Name); err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	return universities, rows.Err()
}

func (repository *PostgresCatalogRepository) FindUniversityByID(ctx context.Context, universityID int64) (models.University, bool, error) {
	var university models.University
	err := repository.pool.QueryRow(ctx, `SELECT id, name FROM university WHERE id = $1`, universityID).Scan(&university.ID, &university.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.University{}, false, nil
	}
	return university, err == nil, err
}

func (repository *PostgresCatalogRepository) UniversityNameExistsInState(ctx context.Context, name string, stateID int64) (bool, error) {
	return queryExists(repository.pool, ctx, `SELECT EXISTS(SELECT 1 FROM university WHERE name = $1 AND state_id = $2)`, name, stateID)
}

func (repository *PostgresCatalogRepository) CreateUniversity(ctx context.Context, university models.University, stateID int64) (models.University, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO university (name, state_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, university.Name, stateID).Scan(&university.ID, &university.Name)
	return university, err
}

func (repository *PostgresCatalogRepository) ListCoursesByUniversity(ctx context.Context, universityID int64) ([]models.Course, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM course WHERE university_id = $1 ORDER BY name ASC`, universityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, rows.Err()
}

func (repository *PostgresCatalogRepository) FindCourseByID(ctx context.Context, courseID int64) (models.Course, bool, error) {
	var course models.Course
	err := repository.pool.QueryRow(ctx, `SELECT id, name FROM course WHERE id = $1`, courseID).Scan(&course.ID, &course.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Course{}, false, nil
	}
	return course, err == nil, err
}

func (repository *PostgresCatalogRepository) CourseNameExistsInUniversity(ctx context.Context, name string, universityID int64) (bool, error) {
	return queryExists(repository.pool, ctx, `SELECT EXISTS(SELECT 1 FROM course WHERE name = $1 AND university_id = $2)`, name, universityID)
}

func (repository *PostgresCatalogRepository) CreateCourse(ctx context.Context, course models.Course, universityID int64) (models.Course, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO course (name, university_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, course.Name, universityID).Scan(&course.ID, &course.Name)
	return course, err
}

func (repository *PostgresCatalogRepository) ListSubjectsByCourse(ctx context.Context, courseID int64) ([]models.Subject, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM subject WHERE course_id = $1 ORDER BY name ASC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := []models.Subject{}
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, rows.Err()
}

func (repository *PostgresCatalogRepository) FindSubjectTreeByID(ctx context.Context, subjectID int64) (models.Subject, bool, error) {
	return findSubjectTree(ctx, repository.pool, `WHERE s.id = $1`, subjectID)
}

func (repository *PostgresCatalogRepository) SubjectNameExistsInCourse(ctx context.Context, name string, courseID int64) (bool, error) {
	return queryExists(repository.pool, ctx, `SELECT EXISTS(SELECT 1 FROM subject WHERE name = $1 AND course_id = $2)`, name, courseID)
}

func (repository *PostgresCatalogRepository) CreateSubject(ctx context.Context, subject models.Subject, courseID int64) (models.Subject, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO subject (name, course_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, subject.Name, courseID).Scan(&subject.ID, &subject.Name)
	return subject, err
}

func queryExists(pool *pgxpool.Pool, ctx context.Context, sql string, args ...any) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, sql, args...).Scan(&exists)
	return exists, err
}

func findSubjectTree(ctx context.Context, pool *pgxpool.Pool, filter string, args ...any) (models.Subject, bool, error) {
	var subject models.Subject
	var course models.Course
	var university models.University
	var state models.State
	err := pool.QueryRow(ctx, `
		SELECT
			s.id, s.name,
			c.id, c.name,
			u.id, u.name,
			st.id, st.code, st.name
		FROM subject s
		JOIN course c ON c.id = s.course_id
		JOIN university u ON u.id = c.university_id
		JOIN state st ON st.id = u.state_id
		`+filter,
		args...,
	).Scan(
		&subject.ID, &subject.Name,
		&course.ID, &course.Name,
		&university.ID, &university.Name,
		&state.ID, &state.Code, &state.Name,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Subject{}, false, nil
	}
	if err != nil {
		return models.Subject{}, false, err
	}
	university.State = &state
	course.University = &university
	subject.Course = &course
	return subject, true, nil
}
