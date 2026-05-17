package examcontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
)

type ApprovedExamFilter struct {
	SubjectID *int64
	ExamYear  *int
	Semester  *int
}

type ExamRepositoryContract interface {
	ListApproved(ctx context.Context, filter ApprovedExamFilter) ([]exam.Exam, error)
	ListPendingBySubject(ctx context.Context, subjectID int64) ([]exam.Exam, error)
	ListPendingByUploader(ctx context.Context, userID int64) ([]exam.Exam, error)
	CountPendingByUploader(ctx context.Context, userID int64) (int64, error)
	FileHashExists(ctx context.Context, fileHash string) (bool, error)
	Create(ctx context.Context, exam exam.Exam) (exam.Exam, error)
	FindByID(ctx context.Context, examID int64) (exam.Exam, bool, error)
	UpdateReview(ctx context.Context, exam exam.Exam) (exam.Exam, error)
}
