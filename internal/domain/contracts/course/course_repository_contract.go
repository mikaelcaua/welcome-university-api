package coursecontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
)

type CourseRepositoryContract interface {
	ListByUniversity(ctx context.Context, universityID int64) ([]course.Course, error)
	FindByID(ctx context.Context, courseID int64) (course.Course, bool, error)
	NameExistsInUniversity(ctx context.Context, name string, universityID int64) (bool, error)
	Create(ctx context.Context, course course.Course, universityID int64) (course.Course, error)
}
