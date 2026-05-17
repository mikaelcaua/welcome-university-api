package examusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
)

type ListMyPendingUseCase struct{ exams examcontract.ExamRepositoryContract }

func NewListMyPendingUseCase(exams examcontract.ExamRepositoryContract) *ListMyPendingUseCase {
	return &ListMyPendingUseCase{exams: exams}
}
func (useCase *ListMyPendingUseCase) Execute(ctx context.Context, currentUser user.User) ([]exam.Exam, error) {
	return useCase.exams.ListPendingByUploader(ctx, currentUser.ID)
}
