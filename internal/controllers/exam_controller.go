package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
	"github.com/mikaelcaua/welcome-university-api/internal/services"
)

type ExamController struct {
	examService     *services.ExamService
	maxUploadSizeMB int64
}

func NewExamController(examService *services.ExamService, maxUploadSizeMB int64) *ExamController {
	return &ExamController{examService: examService, maxUploadSizeMB: maxUploadSizeMB}
}

func (controller *ExamController) ListBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := parsePathInt64(w, r, "subjectId")
	if !ok {
		return
	}
	controller.listApproved(w, r, &subjectID)
}

func (controller *ExamController) ListAll(w http.ResponseWriter, r *http.Request) {
	var subjectID *int64
	rawSubjectID := r.URL.Query().Get("subjectId")
	if rawSubjectID != "" {
		parsedSubjectID, err := strconv.ParseInt(rawSubjectID, 10, 64)
		if err != nil {
			httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "subjectId invalido."))
			return
		}
		subjectID = &parsedSubjectID
	}
	controller.listApproved(w, r, subjectID)
}

func (controller *ExamController) ListPending(w http.ResponseWriter, r *http.Request) {
	stateID, ok := parseRequiredQueryInt64(w, r, "stateId")
	if !ok {
		return
	}
	universityID, ok := parseRequiredQueryInt64(w, r, "universityId")
	if !ok {
		return
	}
	courseID, ok := parseRequiredQueryInt64(w, r, "courseId")
	if !ok {
		return
	}
	subjectID, ok := parseRequiredQueryInt64(w, r, "subjectId")
	if !ok {
		return
	}
	response, err := controller.examService.ListPending(r.Context(), stateID, universityID, courseID, subjectID)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}

func (controller *ExamController) ListMyPending(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r.Context())
	if !ok {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	response, err := controller.examService.ListMyPending(r.Context(), currentUser)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}

func (controller *ExamController) Upload(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r.Context())
	if !ok {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, controller.maxUploadSizeMB*1024*1024)
	if err := r.ParseMultipartForm(controller.maxUploadSizeMB * 1024 * 1024); err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Formulario multipart invalido."))
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Arquivo obrigatorio."))
		return
	}
	defer file.Close()

	payload, err := services.NewUploadPayload(fileHeader.Filename, fileHeader.Header.Get("Content-Type"), file)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	examYear, err := strconv.Atoi(r.FormValue("examYear"))
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Ano da prova invalido."))
		return
	}
	semester, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Semestre invalido."))
		return
	}
	subjectID, err := strconv.ParseInt(r.FormValue("subjectId"), 10, 64)
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "subjectId invalido."))
		return
	}
	periodUnidentified, _ := strconv.ParseBool(r.FormValue("periodUnidentified"))

	response, err := controller.examService.Upload(r.Context(), currentUser, dto.UploadExamRequest{
		ExamYear:           examYear,
		Semester:           semester,
		Type:               models.ExamType(r.FormValue("type")),
		SubjectID:          subjectID,
		PeriodUnidentified: periodUnidentified,
		OriginalFilename:   payload.OriginalFilename,
		ContentType:         payload.ContentType,
		FileBytes:           payload.Bytes,
	})
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, response)
}

func (controller *ExamController) Review(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r.Context())
	if !ok {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	examID, err := strconv.ParseInt(chi.URLParam(r, "examId"), 10, 64)
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return
	}
	var request dto.ReviewExamRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	response, err := controller.examService.Review(r.Context(), currentUser, examID, request)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}

func (controller *ExamController) listApproved(w http.ResponseWriter, r *http.Request, subjectID *int64) {
	response, err := controller.examService.ListApproved(r.Context(), subjectID, r.URL.Query().Get("period"))
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}

func parseRequiredQueryInt64(w http.ResponseWriter, r *http.Request, parameterName string) (int64, bool) {
	value, err := strconv.ParseInt(r.URL.Query().Get(parameterName), 10, 64)
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, parameterName+" invalido."))
		return 0, false
	}
	return value, true
}
