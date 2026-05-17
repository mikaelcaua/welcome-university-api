package subjectusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/subject"
)

type ListSubjectsUseCase struct{ subjects subjectcontract.SubjectRepositoryContract }

func NewListSubjectsUseCase(subjects subjectcontract.SubjectRepositoryContract) *ListSubjectsUseCase {
	return &ListSubjectsUseCase{subjects: subjects}
}
func (useCase *ListSubjectsUseCase) Execute(ctx context.Context, courseID int64) ([]subject.Subject, error) {
	return useCase.subjects.ListByCourse(ctx, courseID)
}
