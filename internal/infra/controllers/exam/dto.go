package examcontroller

import (
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
)

type ReviewExamRequest struct {
	Status     exam.ExamStatus `json:"status"`
	ReviewNote *string             `json:"reviewNote"`
}

type ExamResponse struct {
	ID          int64                               `json:"id"`
	Name        string                              `json:"name"`
	ExamYear    int                                 `json:"examYear"`
	Semester    int                                 `json:"semester"`
	PeriodLabel *string                             `json:"periodLabel"`
	Type        exam.ExamType                   `json:"type"`
	PDFURL      string                              `json:"pdfUrl"`
	Status      exam.ExamStatus                 `json:"status"`
	SubjectID   *int64                              `json:"subjectId"`
	SubjectName *string                             `json:"subjectName"`
	UploadedBy  *usercontroller.UserSummaryResponse `json:"uploadedBy"`
	ReviewedBy  *usercontroller.UserSummaryResponse `json:"reviewedBy"`
	ReviewNote  *string                             `json:"reviewNote"`
	CreatedAt   time.Time                           `json:"createdAt"`
	ReviewedAt  *time.Time                          `json:"reviewedAt"`
}

func ToResponse(exam exam.Exam) ExamResponse {
	var subjectID *int64
	var subjectName *string
	if exam.Subject != nil {
		subjectID = &exam.Subject.ID
		subjectName = &exam.Subject.Name
	}
	return ExamResponse{
		ID: exam.ID, Name: exam.Name, ExamYear: exam.ExamYear, Semester: exam.Semester, PeriodLabel: exam.PeriodLabel,
		Type: exam.Type, PDFURL: exam.PDFURL, Status: exam.Status, SubjectID: subjectID, SubjectName: subjectName,
		UploadedBy: usercontroller.ToSummaryResponse(exam.UploadedBy), ReviewedBy: usercontroller.ToSummaryResponse(exam.ReviewedBy),
		ReviewNote: exam.ReviewNote, CreatedAt: exam.CreatedAt, ReviewedAt: exam.ReviewedAt,
	}
}

func ToResponses(exams []exam.Exam) []ExamResponse {
	response := make([]ExamResponse, 0, len(exams))
	for _, item := range exams {
		response = append(response, ToResponse(item))
	}
	return response
}
