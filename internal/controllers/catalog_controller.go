package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/services"
)

type CatalogController struct {
	catalogService *services.CatalogService
}

func NewCatalogController(catalogService *services.CatalogService) *CatalogController {
	return &CatalogController{catalogService: catalogService}
}

func (controller *CatalogController) ListStates(w http.ResponseWriter, r *http.Request) {
	states, err := controller.catalogService.ListStates(r.Context())
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, states)
}

func (controller *CatalogController) GetStateByCode(w http.ResponseWriter, r *http.Request) {
	state, err := controller.catalogService.GetStateByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, state)
}

func (controller *CatalogController) CreateState(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateStateRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	state, err := controller.catalogService.CreateState(r.Context(), request.Code, request.Name)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, state)
}

func (controller *CatalogController) ListUniversities(w http.ResponseWriter, r *http.Request) {
	stateID, ok := parsePathInt64(w, r, "stateId")
	if !ok {
		return
	}
	universities, err := controller.catalogService.ListUniversities(r.Context(), stateID)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, universities)
}

func (controller *CatalogController) CreateUniversity(w http.ResponseWriter, r *http.Request) {
	stateID, ok := parsePathInt64(w, r, "stateId")
	if !ok {
		return
	}
	var request dto.CreateUniversityRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	university, err := controller.catalogService.CreateUniversity(r.Context(), stateID, request.Name)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, university)
}

func (controller *CatalogController) ListCourses(w http.ResponseWriter, r *http.Request) {
	universityID, ok := parsePathInt64(w, r, "universityId")
	if !ok {
		return
	}
	courses, err := controller.catalogService.ListCourses(r.Context(), universityID)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, courses)
}

func (controller *CatalogController) CreateCourse(w http.ResponseWriter, r *http.Request) {
	universityID, ok := parsePathInt64(w, r, "universityId")
	if !ok {
		return
	}
	var request dto.CreateCourseRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	course, err := controller.catalogService.CreateCourse(r.Context(), universityID, request.Name)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, course)
}

func (controller *CatalogController) ListSubjects(w http.ResponseWriter, r *http.Request) {
	courseID, ok := parsePathInt64(w, r, "courseId")
	if !ok {
		return
	}
	subjects, err := controller.catalogService.ListSubjects(r.Context(), courseID)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, subjects)
}

func (controller *CatalogController) CreateSubject(w http.ResponseWriter, r *http.Request) {
	courseID, ok := parsePathInt64(w, r, "courseId")
	if !ok {
		return
	}
	var request dto.CreateSubjectRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	subject, err := controller.catalogService.CreateSubject(r.Context(), courseID, request.Name)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, subject)
}

func parsePathInt64(w http.ResponseWriter, r *http.Request, parameterName string) (int64, bool) {
	value, err := strconv.ParseInt(chi.URLParam(r, parameterName), 10, 64)
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return 0, false
	}
	return value, true
}
