package subjectcontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
)

type SubjectRepositoryContract interface {
	ListByCourse(ctx context.Context, courseID int64) ([]subject.Subject, error)
	FindTreeByID(ctx context.Context, subjectID int64) (subject.Subject, bool, error)
	NameExistsInCourse(ctx context.Context, name string, courseID int64) (bool, error)
	Create(ctx context.Context, subject subject.Subject, courseID int64) (subject.Subject, error)
}
