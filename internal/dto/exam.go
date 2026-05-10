package dto

import (
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type ReviewExamRequest struct {
	Status     models.ExamStatus `json:"status"`
	ReviewNote *string           `json:"reviewNote"`
}

type UploadExamRequest struct {
	ExamYear           int
	Semester           int
	Type               models.ExamType
	SubjectID          int64
	PeriodUnidentified bool
	OriginalFilename   string
	ContentType         string
	FileBytes           []byte
}

type ExamResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	ExamYear    int                  `json:"examYear"`
	Semester    int                  `json:"semester"`
	PeriodLabel *string              `json:"periodLabel"`
	Type        models.ExamType      `json:"type"`
	PDFURL      string               `json:"pdfUrl"`
	Status      models.ExamStatus    `json:"status"`
	SubjectID   *int64               `json:"subjectId"`
	SubjectName *string              `json:"subjectName"`
	UploadedBy  *UserSummaryResponse `json:"uploadedBy"`
	ReviewedBy  *UserSummaryResponse `json:"reviewedBy"`
	ReviewNote  *string              `json:"reviewNote"`
	CreatedAt   time.Time            `json:"createdAt"`
	ReviewedAt  *time.Time           `json:"reviewedAt"`
}

func ExamToResponse(exam models.Exam) ExamResponse {
	var subjectID *int64
	var subjectName *string
	if exam.Subject != nil {
		subjectID = &exam.Subject.ID
		subjectName = &exam.Subject.Name
	}

	return ExamResponse{
		ID:          exam.ID,
		Name:        exam.Name,
		ExamYear:    exam.ExamYear,
		Semester:    exam.Semester,
		PeriodLabel: exam.PeriodLabel,
		Type:        exam.Type,
		PDFURL:      exam.PDFURL,
		Status:      exam.Status,
		SubjectID:   subjectID,
		SubjectName: subjectName,
		UploadedBy:  UserToSummaryResponse(exam.UploadedBy),
		ReviewedBy:  UserToSummaryResponse(exam.ReviewedBy),
		ReviewNote:  exam.ReviewNote,
		CreatedAt:   exam.CreatedAt,
		ReviewedAt:  exam.ReviewedAt,
	}
}
