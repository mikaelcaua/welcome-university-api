package services

import (
	"context"
	"net/http"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type CatalogRepository interface {
	ListStates(ctx context.Context) ([]models.State, error)
	FindStateByID(ctx context.Context, stateID int64) (models.State, bool, error)
	FindStateByCode(ctx context.Context, code string) (models.State, bool, error)
	CreateState(ctx context.Context, state models.State) (models.State, error)
	ListUniversitiesByState(ctx context.Context, stateID int64) ([]models.University, error)
	FindUniversityByID(ctx context.Context, universityID int64) (models.University, bool, error)
	UniversityNameExistsInState(ctx context.Context, name string, stateID int64) (bool, error)
	CreateUniversity(ctx context.Context, university models.University, stateID int64) (models.University, error)
	ListCoursesByUniversity(ctx context.Context, universityID int64) ([]models.Course, error)
	FindCourseByID(ctx context.Context, courseID int64) (models.Course, bool, error)
	CourseNameExistsInUniversity(ctx context.Context, name string, universityID int64) (bool, error)
	CreateCourse(ctx context.Context, course models.Course, universityID int64) (models.Course, error)
	ListSubjectsByCourse(ctx context.Context, courseID int64) ([]models.Subject, error)
	FindSubjectTreeByID(ctx context.Context, subjectID int64) (models.Subject, bool, error)
	SubjectNameExistsInCourse(ctx context.Context, name string, courseID int64) (bool, error)
	CreateSubject(ctx context.Context, subject models.Subject, courseID int64) (models.Subject, error)
}

type CatalogService struct {
	catalog CatalogRepository
}

func NewCatalogService(catalog CatalogRepository) *CatalogService {
	return &CatalogService{catalog: catalog}
}

func (service *CatalogService) ListStates(ctx context.Context) ([]models.State, error) {
	return service.catalog.ListStates(ctx)
}

func (service *CatalogService) GetStateByCode(ctx context.Context, code string) (models.State, error) {
	state, found, err := service.catalog.FindStateByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		return models.State{}, err
	}
	if !found {
		return models.State{}, httpx.NewHTTPError(http.StatusNotFound, "Estado não encontrado: "+code)
	}
	return state, nil
}

func (service *CatalogService) CreateState(ctx context.Context, code string, name string) (models.State, error) {
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))
	if len(normalizedCode) != 2 || strings.TrimSpace(name) == "" {
		return models.State{}, httpx.NewHTTPError(http.StatusBadRequest, "Sigla e nome sao obrigatorios.")
	}
	_, found, err := service.catalog.FindStateByCode(ctx, normalizedCode)
	if err != nil {
		return models.State{}, err
	}
	if found {
		return models.State{}, httpx.NewHTTPError(http.StatusConflict, "Estado ja cadastrado: "+normalizedCode)
	}
	return service.catalog.CreateState(ctx, models.State{Code: normalizedCode, Name: strings.TrimSpace(name)})
}

func (service *CatalogService) ListUniversities(ctx context.Context, stateID int64) ([]models.University, error) {
	return service.catalog.ListUniversitiesByState(ctx, stateID)
}

func (service *CatalogService) CreateUniversity(ctx context.Context, stateID int64, name string) (models.University, error) {
	_, found, err := service.catalog.FindStateByID(ctx, stateID)
	if err != nil {
		return models.University{}, err
	}
	if !found {
		return models.University{}, httpx.NewHTTPError(http.StatusNotFound, "Estado nao encontrado.")
	}
	normalizedName := strings.TrimSpace(name)
	exists, err := service.catalog.UniversityNameExistsInState(ctx, normalizedName, stateID)
	if err != nil {
		return models.University{}, err
	}
	if exists {
		return models.University{}, httpx.NewHTTPError(http.StatusConflict, "Universidade ja cadastrada neste estado.")
	}
	return service.catalog.CreateUniversity(ctx, models.University{Name: normalizedName}, stateID)
}

func (service *CatalogService) ListCourses(ctx context.Context, universityID int64) ([]models.Course, error) {
	return service.catalog.ListCoursesByUniversity(ctx, universityID)
}

func (service *CatalogService) CreateCourse(ctx context.Context, universityID int64, name string) (models.Course, error) {
	_, found, err := service.catalog.FindUniversityByID(ctx, universityID)
	if err != nil {
		return models.Course{}, err
	}
	if !found {
		return models.Course{}, httpx.NewHTTPError(http.StatusNotFound, "Universidade nao encontrada.")
	}
	normalizedName := strings.TrimSpace(name)
	exists, err := service.catalog.CourseNameExistsInUniversity(ctx, normalizedName, universityID)
	if err != nil {
		return models.Course{}, err
	}
	if exists {
		return models.Course{}, httpx.NewHTTPError(http.StatusConflict, "Curso ja cadastrado nesta universidade.")
	}
	return service.catalog.CreateCourse(ctx, models.Course{Name: normalizedName}, universityID)
}

func (service *CatalogService) ListSubjects(ctx context.Context, courseID int64) ([]models.Subject, error) {
	return service.catalog.ListSubjectsByCourse(ctx, courseID)
}

func (service *CatalogService) CreateSubject(ctx context.Context, courseID int64, name string) (models.Subject, error) {
	_, found, err := service.catalog.FindCourseByID(ctx, courseID)
	if err != nil {
		return models.Subject{}, err
	}
	if !found {
		return models.Subject{}, httpx.NewHTTPError(http.StatusNotFound, "Curso nao encontrado.")
	}
	normalizedName := strings.TrimSpace(name)
	exists, err := service.catalog.SubjectNameExistsInCourse(ctx, normalizedName, courseID)
	if err != nil {
		return models.Subject{}, err
	}
	if exists {
		return models.Subject{}, httpx.NewHTTPError(http.StatusConflict, "Disciplina ja cadastrada neste curso.")
	}
	return service.catalog.CreateSubject(ctx, models.Subject{Name: normalizedName}, courseID)
}
