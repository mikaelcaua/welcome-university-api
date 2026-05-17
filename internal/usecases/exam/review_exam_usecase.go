package examusecase

import (
	"context"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type ReviewExamUseCase struct {
	exams   examcontract.ExamRepositoryContract
	storage storagecontract.ObjectStorageContract
}

func NewReviewExamUseCase(exams examcontract.ExamRepositoryContract, storage storagecontract.ObjectStorageContract) *ReviewExamUseCase {
	return &ReviewExamUseCase{exams: exams, storage: storage}
}
func (useCase *ReviewExamUseCase) Execute(ctx context.Context, currentUser user.User, examID int64, input ReviewInput) (exam.Exam, error) {
	if !input.Status.IsReviewResult() {
		return exam.Exam{}, domainerrors.Validation("Status de revisao invalido.")
	}
	exam, found, err := useCase.exams.FindByID(ctx, examID)
	if err != nil {
		return exam.Exam{}, err
	}
	if !found {
		return exam.Exam{}, domainerrors.NotFound("Prova nao encontrada.")
	}
	if !exam.CanBeReviewed() {
		return exam.Exam{}, domainerrors.Conflict("Apenas provas pendentes podem ser validadas.")
	}
	if input.Status == exam.ExamStatusRejected {
		if err := useCase.storage.DeleteObjectByKey(ctx, exam.StorageKey); err != nil {
			return exam.Exam{}, err
		}
		exam.PDFURL = ""
		exam.StorageKey = nil
	}
	now := time.Now()
	exam.Status = input.Status
	exam.ReviewedBy = &currentUser
	exam.ReviewedAt = &now
	exam.ReviewNote = normalizeOptionalText(input.ReviewNote)
	return useCase.exams.UpdateReview(ctx, exam)
}
