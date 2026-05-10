package services

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
	"github.com/mikaelcaua/welcome-university-api/internal/repositories"
)

func TestUploadByAdminIsAutomaticallyApproved(t *testing.T) {
	examRepository := newFakeExamRepository()
	catalogRepository := fakeCatalogRepository{subject: completeSubjectTree()}
	storage := &fakeObjectStorage{}
	service := NewExamService(examRepository, catalogRepository, storage, passthroughOptimizer{})

	admin := models.User{ID: 10, Name: "Admin", Email: "admin@example.com", Role: models.RoleAdmin}
	response, err := service.Upload(context.Background(), admin, dto.UploadExamRequest{
		ExamYear:         2024,
		Semester:         1,
		Type:             models.ExamTypeFirst,
		SubjectID:        30,
		OriginalFilename: "prova.pdf",
		ContentType:      "application/pdf",
		FileBytes:        []byte("pdf-content"),
	})
	if err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}
	if response.Status != models.ExamStatusApproved {
		t.Fatalf("expected APPROVED status, got %s", response.Status)
	}
	if response.ReviewedBy == nil || response.ReviewedBy.ID != admin.ID {
		t.Fatal("expected admin reviewer in response")
	}
	if storage.uploadedSubjectID != 30 {
		t.Fatalf("expected upload under subject 30, got %d", storage.uploadedSubjectID)
	}
}

func TestUploadRejectsRegularUserOverPendingLimit(t *testing.T) {
	examRepository := newFakeExamRepository()
	examRepository.pendingCount = maxPendingExamsForRegularUser
	service := NewExamService(examRepository, fakeCatalogRepository{subject: completeSubjectTree()}, &fakeObjectStorage{}, passthroughOptimizer{})

	regularUser := models.User{ID: 11, Role: models.RoleUser}
	_, err := service.Upload(context.Background(), regularUser, dto.UploadExamRequest{
		ExamYear:         2024,
		Semester:         1,
		Type:             models.ExamTypeFirst,
		SubjectID:        30,
		OriginalFilename: "prova.pdf",
		ContentType:      "application/pdf",
		FileBytes:        []byte("pdf-content"),
	})
	if err == nil {
		t.Fatal("expected pending limit error")
	}
	var httpError httpx.HTTPError
	if !errors.As(err, &httpError) || httpError.StatusCode != http.StatusConflict {
		t.Fatalf("expected conflict HTTP error, got %v", err)
	}
}

type fakeExamRepository struct {
	exams        map[int64]models.Exam
	nextID       int64
	pendingCount int64
	hashes       map[string]bool
}

func newFakeExamRepository() *fakeExamRepository {
	return &fakeExamRepository{
		exams:  map[int64]models.Exam{},
		nextID: 1,
		hashes: map[string]bool{},
	}
}

func (repository *fakeExamRepository) ListApproved(ctx context.Context, filter repositories.ApprovedExamFilter) ([]models.Exam, error) {
	return nil, nil
}

func (repository *fakeExamRepository) ListPendingBySubject(ctx context.Context, subjectID int64) ([]models.Exam, error) {
	return nil, nil
}

func (repository *fakeExamRepository) ListPendingByUploader(ctx context.Context, userID int64) ([]models.Exam, error) {
	return nil, nil
}

func (repository *fakeExamRepository) CountPendingByUploader(ctx context.Context, userID int64) (int64, error) {
	return repository.pendingCount, nil
}

func (repository *fakeExamRepository) FileHashExists(ctx context.Context, fileHash string) (bool, error) {
	return repository.hashes[fileHash], nil
}

func (repository *fakeExamRepository) Create(ctx context.Context, exam models.Exam) (models.Exam, error) {
	exam.ID = repository.nextID
	repository.nextID++
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = exam.CreatedAt
	if exam.FileHash != nil {
		repository.hashes[*exam.FileHash] = true
	}
	repository.exams[exam.ID] = exam
	return exam, nil
}

func (repository *fakeExamRepository) FindByID(ctx context.Context, examID int64) (models.Exam, bool, error) {
	exam, found := repository.exams[examID]
	return exam, found, nil
}

func (repository *fakeExamRepository) UpdateReview(ctx context.Context, exam models.Exam) (models.Exam, error) {
	repository.exams[exam.ID] = exam
	return exam, nil
}

type fakeCatalogRepository struct {
	subject models.Subject
}

func (repository fakeCatalogRepository) ListStates(ctx context.Context) ([]models.State, error) { return nil, nil }
func (repository fakeCatalogRepository) FindStateByID(ctx context.Context, stateID int64) (models.State, bool, error) {
	return models.State{}, true, nil
}
func (repository fakeCatalogRepository) FindStateByCode(ctx context.Context, code string) (models.State, bool, error) {
	return models.State{}, true, nil
}
func (repository fakeCatalogRepository) CreateState(ctx context.Context, state models.State) (models.State, error) {
	return state, nil
}
func (repository fakeCatalogRepository) ListUniversitiesByState(ctx context.Context, stateID int64) ([]models.University, error) {
	return nil, nil
}
func (repository fakeCatalogRepository) FindUniversityByID(ctx context.Context, universityID int64) (models.University, bool, error) {
	return models.University{}, true, nil
}
func (repository fakeCatalogRepository) UniversityNameExistsInState(ctx context.Context, name string, stateID int64) (bool, error) {
	return false, nil
}
func (repository fakeCatalogRepository) CreateUniversity(ctx context.Context, university models.University, stateID int64) (models.University, error) {
	return university, nil
}
func (repository fakeCatalogRepository) ListCoursesByUniversity(ctx context.Context, universityID int64) ([]models.Course, error) {
	return nil, nil
}
func (repository fakeCatalogRepository) FindCourseByID(ctx context.Context, courseID int64) (models.Course, bool, error) {
	return models.Course{}, true, nil
}
func (repository fakeCatalogRepository) CourseNameExistsInUniversity(ctx context.Context, name string, universityID int64) (bool, error) {
	return false, nil
}
func (repository fakeCatalogRepository) CreateCourse(ctx context.Context, course models.Course, universityID int64) (models.Course, error) {
	return course, nil
}
func (repository fakeCatalogRepository) ListSubjectsByCourse(ctx context.Context, courseID int64) ([]models.Subject, error) {
	return nil, nil
}
func (repository fakeCatalogRepository) FindSubjectTreeByID(ctx context.Context, subjectID int64) (models.Subject, bool, error) {
	return repository.subject, repository.subject.ID == subjectID, nil
}
func (repository fakeCatalogRepository) SubjectNameExistsInCourse(ctx context.Context, name string, courseID int64) (bool, error) {
	return false, nil
}
func (repository fakeCatalogRepository) CreateSubject(ctx context.Context, subject models.Subject, courseID int64) (models.Subject, error) {
	return subject, nil
}

type fakeObjectStorage struct {
	uploadedSubjectID int64
	deletedStorageKey *string
}

func (storage *fakeObjectStorage) UploadExam(ctx context.Context, payload UploadPayload, subjectID int64) (StoredObject, error) {
	storage.uploadedSubjectID = subjectID
	return StoredObject{Key: "subjects/30/file.pdf", URL: "https://cdn.example/exams-bucket/subjects/30/file.pdf"}, nil
}

func (storage *fakeObjectStorage) DeleteObjectByKey(ctx context.Context, storageKey *string) error {
	storage.deletedStorageKey = storageKey
	return nil
}

type passthroughOptimizer struct{}

func (passthroughOptimizer) Optimize(payload UploadPayload) UploadPayload {
	return payload
}

func completeSubjectTree() models.Subject {
	state := models.State{ID: 1, Code: "MA", Name: "Maranhao"}
	university := models.University{ID: 2, Name: "UFMA", State: &state}
	course := models.Course{ID: 3, Name: "Ciencia da Computacao", University: &university}
	return models.Subject{ID: 30, Name: "Algoritmos", Course: &course}
}
