package exam

import (
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
)

type ExamType string

const (
	ExamTypeFirst    ExamType = "PROVA1"
	ExamTypeSecond   ExamType = "PROVA2"
	ExamTypeThird    ExamType = "PROVA3"
	ExamTypeRecovery ExamType = "RECUPERACAO"
	ExamTypeFinal    ExamType = "FINAL"
)

func (examType ExamType) IsValid() bool {
	switch examType {
	case ExamTypeFirst, ExamTypeSecond, ExamTypeThird, ExamTypeRecovery, ExamTypeFinal:
		return true
	default:
		return false
	}
}

type ExamStatus string

const (
	ExamStatusPending  ExamStatus = "PENDING"
	ExamStatusApproved ExamStatus = "APPROVED"
	ExamStatusRejected ExamStatus = "REJECTED"
)

func (status ExamStatus) IsReviewResult() bool {
	return status == ExamStatusApproved || status == ExamStatusRejected
}

type Exam struct {
	ID          int64
	Name        string
	ExamYear    int
	Semester    int
	PeriodLabel *string
	Type        ExamType
	PDFURL      string
	StorageKey  *string
	FileHash    *string
	Status      ExamStatus
	Subject     *subject.Subject
	UploadedBy  *user.User
	ReviewedBy  *user.User
	ReviewNote  *string
	CreatedAt   time.Time
	ReviewedAt  *time.Time
	UpdatedAt   time.Time
}

func (exam Exam) CanBeReviewed() bool {
	return exam.Status == ExamStatusPending
}
