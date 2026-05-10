package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
	"github.com/mikaelcaua/welcome-university-api/internal/repositories"
)

const maxPendingExamsForRegularUser = 5
const unidentifiedPeriodLabel = "PERIODO_NAO_IDENTIFICADO"

var supportedUploadExtensions = map[string]bool{
	"pdf": true, "png": true, "jpg": true, "jpeg": true, "webp": true, "gif": true,
}

type ExamRepository interface {
	ListApproved(ctx context.Context, filter repositories.ApprovedExamFilter) ([]models.Exam, error)
	ListPendingBySubject(ctx context.Context, subjectID int64) ([]models.Exam, error)
	ListPendingByUploader(ctx context.Context, userID int64) ([]models.Exam, error)
	CountPendingByUploader(ctx context.Context, userID int64) (int64, error)
	FileHashExists(ctx context.Context, fileHash string) (bool, error)
	Create(ctx context.Context, exam models.Exam) (models.Exam, error)
	FindByID(ctx context.Context, examID int64) (models.Exam, bool, error)
	UpdateReview(ctx context.Context, exam models.Exam) (models.Exam, error)
}

type UploadPayloadOptimizer interface {
	Optimize(payload UploadPayload) UploadPayload
}

type ExamService struct {
	exams    ExamRepository
	catalog  CatalogRepository
	storage  ObjectStorage
	optimizer UploadPayloadOptimizer
}

func NewExamService(exams ExamRepository, catalog CatalogRepository, storage ObjectStorage, optimizer UploadPayloadOptimizer) *ExamService {
	return &ExamService{exams: exams, catalog: catalog, storage: storage, optimizer: optimizer}
}

func (service *ExamService) ListApproved(ctx context.Context, subjectID *int64, period string) ([]dto.ExamResponse, error) {
	if strings.TrimSpace(period) != "" && subjectID == nil {
		return nil, httpx.NewHTTPError(http.StatusBadRequest, "O filtro period exige subjectId.")
	}
	filter := repositories.ApprovedExamFilter{SubjectID: subjectID}
	if strings.TrimSpace(period) != "" {
		examYear, semester, err := parsePeriod(period)
		if err != nil {
			return nil, err
		}
		filter.ExamYear = &examYear
		filter.Semester = &semester
	}
	exams, err := service.exams.ListApproved(ctx, filter)
	if err != nil {
		return nil, err
	}
	return examResponses(exams), nil
}

func (service *ExamService) ListPending(ctx context.Context, stateID int64, universityID int64, courseID int64, subjectID int64) ([]dto.ExamResponse, error) {
	subject, found, err := service.catalog.FindSubjectTreeByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, httpx.NewHTTPError(http.StatusNotFound, "Materia nao encontrada.")
	}
	if subject.Course == nil || subject.Course.ID != courseID {
		return nil, httpx.NewHTTPError(http.StatusBadRequest, "subjectId nao pertence ao courseId informado.")
	}
	if subject.Course.University == nil || subject.Course.University.ID != universityID {
		return nil, httpx.NewHTTPError(http.StatusBadRequest, "courseId nao pertence ao universityId informado.")
	}
	if subject.Course.University.State == nil || subject.Course.University.State.ID != stateID {
		return nil, httpx.NewHTTPError(http.StatusBadRequest, "universityId nao pertence ao stateId informado.")
	}
	exams, err := service.exams.ListPendingBySubject(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	return examResponses(exams), nil
}

func (service *ExamService) ListMyPending(ctx context.Context, currentUser models.User) ([]dto.ExamResponse, error) {
	exams, err := service.exams.ListPendingByUploader(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	return examResponses(exams), nil
}

func (service *ExamService) Upload(ctx context.Context, currentUser models.User, request dto.UploadExamRequest) (dto.ExamResponse, error) {
	if err := validateUploadRequest(request); err != nil {
		return dto.ExamResponse{}, err
	}
	if currentUser.Role != models.RoleAdmin && currentUser.Role != models.RoleDev {
		pendingCount, err := service.exams.CountPendingByUploader(ctx, currentUser.ID)
		if err != nil {
			return dto.ExamResponse{}, err
		}
		if pendingCount >= maxPendingExamsForRegularUser {
			return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusConflict, "Somente ADMIN/DEV podem ultrapassar o limite de 5 provas pendentes.")
		}
	}
	subject, found, err := service.catalog.FindSubjectTreeByID(ctx, request.SubjectID)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	if !found {
		return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusNotFound, "Materia nao encontrada.")
	}

	payload := service.optimizer.Optimize(UploadPayload{
		OriginalFilename: request.OriginalFilename,
		ContentType:      request.ContentType,
		Bytes:            request.FileBytes,
	})
	fileHash := sha256Hex(payload.Bytes)
	exists, err := service.exams.FileHashExists(ctx, fileHash)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	if exists {
		return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusConflict, "Arquivo ja existe na base.")
	}
	storedObject, err := service.storage.UploadExam(ctx, payload, subject.ID)
	if err != nil {
		return dto.ExamResponse{}, err
	}

	now := time.Now()
	fileHashPointer := fileHash
	storageKeyPointer := storedObject.Key
	exam := models.Exam{
		Name:       buildExamName(subject, request),
		ExamYear:   request.ExamYear,
		Semester:   request.Semester,
		Type:       request.Type,
		PDFURL:     storedObject.URL,
		StorageKey: &storageKeyPointer,
		FileHash:   &fileHashPointer,
		Status:     models.ExamStatusPending,
		Subject:    &subject,
		UploadedBy: &currentUser,
	}
	if request.PeriodUnidentified {
		periodLabel := unidentifiedPeriodLabel
		exam.PeriodLabel = &periodLabel
	}
	if currentUser.Role == models.RoleAdmin || currentUser.Role == models.RoleDev {
		reviewNote := "Aprovacao automatica: upload por ADMIN/DEV."
		exam.Status = models.ExamStatusApproved
		exam.ReviewedBy = &currentUser
		exam.ReviewedAt = &now
		exam.ReviewNote = &reviewNote
	}
	createdExam, err := service.exams.Create(ctx, exam)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	createdExam.Subject = &subject
	createdExam.UploadedBy = &currentUser
	createdExam.ReviewedBy = exam.ReviewedBy
	createdExam.ReviewedAt = exam.ReviewedAt
	createdExam.ReviewNote = exam.ReviewNote
	return dto.ExamToResponse(createdExam), nil
}

func (service *ExamService) Review(ctx context.Context, currentUser models.User, examID int64, request dto.ReviewExamRequest) (dto.ExamResponse, error) {
	if !request.Status.IsReviewResult() {
		return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusBadRequest, "Status de revisao invalido.")
	}
	exam, found, err := service.exams.FindByID(ctx, examID)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	if !found {
		return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusNotFound, "Prova nao encontrada.")
	}
	if exam.Status != models.ExamStatusPending {
		return dto.ExamResponse{}, httpx.NewHTTPError(http.StatusConflict, "Apenas provas pendentes podem ser validadas.")
	}
	if request.Status == models.ExamStatusRejected {
		if err := service.storage.DeleteObjectByKey(ctx, exam.StorageKey); err != nil {
			return dto.ExamResponse{}, err
		}
		exam.PDFURL = ""
		exam.StorageKey = nil
	}
	now := time.Now()
	exam.Status = request.Status
	exam.ReviewedBy = &currentUser
	exam.ReviewedAt = &now
	exam.ReviewNote = normalizeOptionalText(request.ReviewNote)
	exam, err = service.exams.UpdateReview(ctx, exam)
	if err != nil {
		return dto.ExamResponse{}, err
	}
	return dto.ExamToResponse(exam), nil
}

func examResponses(exams []models.Exam) []dto.ExamResponse {
	responses := make([]dto.ExamResponse, 0, len(exams))
	for _, exam := range exams {
		responses = append(responses, dto.ExamToResponse(exam))
	}
	return responses
}

func parsePeriod(period string) (int, int, error) {
	parts := strings.Split(period, ".")
	if len(parts) != 2 {
		return 0, 0, httpx.NewHTTPError(http.StatusBadRequest, "Period deve seguir o formato AAAA.S.")
	}
	examYear, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, httpx.NewHTTPError(http.StatusBadRequest, "Period deve conter apenas numeros.")
	}
	semester, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, httpx.NewHTTPError(http.StatusBadRequest, "Period deve conter apenas numeros.")
	}
	if semester < 1 || semester > 2 {
		return 0, 0, httpx.NewHTTPError(http.StatusBadRequest, "Semestre invalido.")
	}
	return examYear, semester, nil
}

func validateUploadRequest(request dto.UploadExamRequest) error {
	if request.ExamYear < 2000 {
		return httpx.NewHTTPError(http.StatusBadRequest, "Ano da prova invalido.")
	}
	if request.Semester < 1 || request.Semester > 2 {
		return httpx.NewHTTPError(http.StatusBadRequest, "Semestre invalido.")
	}
	if !request.Type.IsValid() {
		return httpx.NewHTTPError(http.StatusBadRequest, "Tipo de prova invalido.")
	}
	if len(request.FileBytes) == 0 {
		return httpx.NewHTTPError(http.StatusBadRequest, "Arquivo obrigatorio.")
	}
	lowerFilename := strings.ToLower(request.OriginalFilename)
	lowerContentType := strings.ToLower(request.ContentType)
	extension := strings.TrimPrefix(strings.ToLower(strings.TrimPrefix(fileExtension(lowerFilename), ".")), ".")
	validFile := lowerContentType == "application/pdf" || strings.HasPrefix(lowerContentType, "image/") || supportedUploadExtensions[extension]
	if !validFile {
		return httpx.NewHTTPError(http.StatusBadRequest, "Somente arquivos PDF ou imagens sao aceitos.")
	}
	return nil
}

func buildExamName(subject models.Subject, request dto.UploadExamRequest) string {
	if request.PeriodUnidentified {
		return fmt.Sprintf("%s - %s - %s", strings.TrimSpace(subject.Name), unidentifiedPeriodLabel, formatExamType(request.Type))
	}
	return fmt.Sprintf("%s - %d.%d - %s", strings.TrimSpace(subject.Name), request.ExamYear, request.Semester, formatExamType(request.Type))
}

func formatExamType(examType models.ExamType) string {
	return strings.ToUpper(strings.ReplaceAll(string(examType), "_", " "))
}

func fileExtension(filename string) string {
	lastDotIndex := strings.LastIndex(filename, ".")
	if lastDotIndex < 0 || lastDotIndex == len(filename)-1 {
		return ""
	}
	return filename[lastDotIndex+1:]
}

func normalizeOptionalText(value *string) *string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	return &normalized
}

func sha256Hex(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
