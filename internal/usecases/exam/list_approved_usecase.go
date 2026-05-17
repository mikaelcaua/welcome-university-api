package examusecase

import (
	"context"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type ListApprovedUseCase struct{ exams examcontract.ExamRepositoryContract }

func NewListApprovedUseCase(exams examcontract.ExamRepositoryContract) *ListApprovedUseCase {
	return &ListApprovedUseCase{exams: exams}
}
func (useCase *ListApprovedUseCase) Execute(ctx context.Context, subjectID *int64, period string) ([]exam.Exam, error) {
	if strings.TrimSpace(period) != "" && subjectID == nil {
		return nil, domainerrors.Validation("O filtro period exige subjectId.")
	}
	filter := examcontract.ApprovedExamFilter{SubjectID: subjectID}
	if strings.TrimSpace(period) != "" {
		examYear, semester, err := parsePeriod(period)
		if err != nil {
			return nil, err
		}
		filter.ExamYear = &examYear
		filter.Semester = &semester
	}
	return useCase.exams.ListApproved(ctx, filter)
}
