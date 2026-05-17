package courseusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/course"
)

type ListCoursesUseCase struct{ courses coursecontract.CourseRepositoryContract }

func NewListCoursesUseCase(courses coursecontract.CourseRepositoryContract) *ListCoursesUseCase {
	return &ListCoursesUseCase{courses: courses}
}
func (useCase *ListCoursesUseCase) Execute(ctx context.Context, universityID int64) ([]course.Course, error) {
	return useCase.courses.ListByUniversity(ctx, universityID)
}
