package examusecase

import (
	"context"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
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
	foundExam, found, err := useCase.exams.FindByID(ctx, examID)
	if err != nil {
		return exam.Exam{}, err
	}
	if !found {
		return exam.Exam{}, domainerrors.NotFound("Prova nao encontrada.")
	}
	if !foundExam.CanBeReviewed() {
		return exam.Exam{}, domainerrors.Conflict("Apenas provas pendentes podem ser validadas.")
	}
	if input.Status == exam.ExamStatusRejected {
		if err := useCase.storage.DeleteObjectByKey(ctx, foundExam.StorageKey); err != nil {
			return exam.Exam{}, err
		}
		foundExam.PDFURL = ""
		foundExam.StorageKey = nil
	}
	now := time.Now()
	foundExam.Status = input.Status
	foundExam.ReviewedBy = &currentUser
	foundExam.ReviewedAt = &now
	foundExam.ReviewNote = normalizeOptionalText(input.ReviewNote)
	return useCase.exams.UpdateReview(ctx, foundExam)
}
