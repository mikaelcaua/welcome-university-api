package courserepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
)

type PostgresCourseRepositoryImpl struct{ pool *pgxpool.Pool }

func NewPostgresCourseRepositoryImpl(pool *pgxpool.Pool) *PostgresCourseRepositoryImpl {
	return &PostgresCourseRepositoryImpl{pool: pool}
}

func (repository *PostgresCourseRepositoryImpl) ListByUniversity(ctx context.Context, universityID int64) ([]course.Course, error) {
	rows, err := repository.pool.Query(ctx, `SELECT id, name FROM course WHERE university_id = $1 ORDER BY name ASC`, universityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []course.Course{}
	for rows.Next() {
		var course course.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, rows.Err()
}

func (repository *PostgresCourseRepositoryImpl) FindByID(ctx context.Context, courseID int64) (course.Course, bool, error) {
	var foundCourse course.Course
	err := repository.pool.QueryRow(ctx, `SELECT id, name FROM course WHERE id = $1`, courseID).Scan(&foundCourse.ID, &foundCourse.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return course.Course{}, false, nil
	}
	return foundCourse, err == nil, err
}

func (repository *PostgresCourseRepositoryImpl) NameExistsInUniversity(ctx context.Context, name string, universityID int64) (bool, error) {
	var exists bool
	err := repository.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM course WHERE name = $1 AND university_id = $2)`, name, universityID).Scan(&exists)
	return exists, err
}

func (repository *PostgresCourseRepositoryImpl) Create(ctx context.Context, course course.Course, universityID int64) (course.Course, error) {
	err := repository.pool.QueryRow(ctx, `
		INSERT INTO course (name, university_id)
		VALUES ($1, $2)
		RETURNING id, name
	`, course.Name, universityID).Scan(&course.ID, &course.Name)
	return course, err
}
