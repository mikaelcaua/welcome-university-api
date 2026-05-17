package examusecase

import (
	"context"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type UploadExamUseCase struct {
	exams     examcontract.ExamRepositoryContract
	subjects  subjectcontract.SubjectRepositoryContract
	storage   storagecontract.ObjectStorageContract
	optimizer UploadPayloadOptimizer
}

func NewUploadExamUseCase(exams examcontract.ExamRepositoryContract, subjects subjectcontract.SubjectRepositoryContract, storage storagecontract.ObjectStorageContract, optimizer UploadPayloadOptimizer) *UploadExamUseCase {
	return &UploadExamUseCase{exams: exams, subjects: subjects, storage: storage, optimizer: optimizer}
}
func (useCase *UploadExamUseCase) Execute(ctx context.Context, currentUser user.User, input UploadInput) (exam.Exam, error) {
	if err := validateUploadInput(input); err != nil {
		return exam.Exam{}, err
	}
	if currentUser.Role != user.RoleAdmin && currentUser.Role != user.RoleDev {
		pendingCount, err := useCase.exams.CountPendingByUploader(ctx, currentUser.ID)
		if err != nil {
			return exam.Exam{}, err
		}
		if pendingCount >= maxPendingExamsForRegularUser {
			return exam.Exam{}, domainerrors.Conflict("Somente ADMIN/DEV podem ultrapassar o limite de 5 provas pendentes.")
		}
	}
	subject, found, err := useCase.subjects.FindTreeByID(ctx, input.SubjectID)
	if err != nil {
		return exam.Exam{}, err
	}
	if !found {
		return exam.Exam{}, domainerrors.NotFound("Materia nao encontrada.")
	}
	payload := useCase.optimizer.Optimize(storagecontract.UploadPayload{OriginalFilename: input.OriginalFilename, ContentType: input.ContentType, Bytes: input.FileBytes})
	fileHash := sha256Hex(payload.Bytes)
	exists, err := useCase.exams.FileHashExists(ctx, fileHash)
	if err != nil {
		return exam.Exam{}, err
	}
	if exists {
		return exam.Exam{}, domainerrors.Conflict("Arquivo ja existe na base.")
	}
	storedObject, err := useCase.storage.UploadExam(ctx, payload, subject.ID)
	if err != nil {
		return exam.Exam{}, err
	}
	now := time.Now()
	fileHashPointer := fileHash
	storageKeyPointer := storedObject.Key
	exam := exam.Exam{Name: buildExamName(subject, input), ExamYear: input.ExamYear, Semester: input.Semester, Type: input.Type, PDFURL: storedObject.URL, StorageKey: &storageKeyPointer, FileHash: &fileHashPointer, Status: exam.ExamStatusPending, Subject: &subject, UploadedBy: &currentUser}
	if input.PeriodUnidentified {
		periodLabel := unidentifiedPeriodLabel
		exam.PeriodLabel = &periodLabel
	}
	if currentUser.Role == user.RoleAdmin || currentUser.Role == user.RoleDev {
		reviewNote := "Aprovacao automatica: upload por ADMIN/DEV."
		exam.Status = exam.ExamStatusApproved
		exam.ReviewedBy = &currentUser
		exam.ReviewedAt = &now
		exam.ReviewNote = &reviewNote
	}
	createdExam, err := useCase.exams.Create(ctx, exam)
	if err != nil {
		return exam.Exam{}, err
	}
	createdExam.Subject = &subject
	createdExam.UploadedBy = &currentUser
	createdExam.ReviewedBy = exam.ReviewedBy
	createdExam.ReviewedAt = exam.ReviewedAt
	createdExam.ReviewNote = exam.ReviewNote
	return createdExam, nil
}
