package examcontroller

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/exam"
)

type ExamController struct {
	listApproved  *examusecase.ListApprovedUseCase
	listPending   *examusecase.ListPendingUseCase
	listMyPending *examusecase.ListMyPendingUseCase
	upload        *examusecase.UploadExamUseCase
	review        *examusecase.ReviewExamUseCase
	maxUploadMB   int64
}

func NewExamController(listApproved *examusecase.ListApprovedUseCase, listPending *examusecase.ListPendingUseCase, listMyPending *examusecase.ListMyPendingUseCase, upload *examusecase.UploadExamUseCase, review *examusecase.ReviewExamUseCase, maxUploadMB int64) *ExamController {
	return &ExamController{listApproved: listApproved, listPending: listPending, listMyPending: listMyPending, upload: upload, review: review, maxUploadMB: maxUploadMB}
}
func (controller *ExamController) ListBySubject(ctx *gin.Context) {
	subjectID, ok := parsePathInt64(ctx, "subjectId")
	if !ok {
		return
	}
	controller.listApprovedResponse(ctx, &subjectID)
}
func (controller *ExamController) ListAll(ctx *gin.Context) {
	var subjectID *int64
	if rawSubjectID := ctx.Query("subjectId"); rawSubjectID != "" {
		parsedSubjectID, err := strconv.ParseInt(rawSubjectID, 10, 64)
		if err != nil {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "subjectId invalido."))
			return
		}
		subjectID = &parsedSubjectID
	}
	controller.listApprovedResponse(ctx, subjectID)
}
func (controller *ExamController) ListPending(ctx *gin.Context) {
	stateID, ok := parseRequiredQueryInt64(ctx, "stateId")
	if !ok {
		return
	}
	universityID, ok := parseRequiredQueryInt64(ctx, "universityId")
	if !ok {
		return
	}
	courseID, ok := parseRequiredQueryInt64(ctx, "courseId")
	if !ok {
		return
	}
	subjectID, ok := parseRequiredQueryInt64(ctx, "subjectId")
	if !ok {
		return
	}
	response, err := controller.listPending.Execute(ctx.Request.Context(), stateID, universityID, courseID, subjectID)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponses(response))
}
func (controller *ExamController) ListMyPending(ctx *gin.Context) {
	currentUser, ok := middleware.CurrentUser(ctx)
	if !ok {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	response, err := controller.listMyPending.Execute(ctx.Request.Context(), currentUser)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponses(response))
}
func (controller *ExamController) Upload(ctx *gin.Context) {
	currentUser, ok := middleware.CurrentUser(ctx)
	if !ok {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, controller.maxUploadMB*1024*1024)
	if err := ctx.Request.ParseMultipartForm(controller.maxUploadMB * 1024 * 1024); err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Formulario multipart invalido."))
		return
	}
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Arquivo obrigatorio."))
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Falha ao ler arquivo enviado."))
		return
	}
	examYear, err := strconv.Atoi(ctx.Request.FormValue("examYear"))
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Ano da prova invalido."))
		return
	}
	semester, err := strconv.Atoi(ctx.Request.FormValue("semester"))
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Semestre invalido."))
		return
	}
	subjectID, err := strconv.ParseInt(ctx.Request.FormValue("subjectId"), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "subjectId invalido."))
		return
	}
	periodUnidentified, _ := strconv.ParseBool(ctx.Request.FormValue("periodUnidentified"))
	response, err := controller.upload.Execute(ctx.Request.Context(), currentUser, examusecase.UploadInput{
		ExamYear: examYear, Semester: semester, Type: exam.ExamType(ctx.Request.FormValue("type")), SubjectID: subjectID,
		PeriodUnidentified: periodUnidentified, OriginalFilename: fileHeader.Filename, ContentType: fileHeader.Header.Get("Content-Type"), FileBytes: fileBytes,
	})
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusCreated, ToResponse(response))
}
func (controller *ExamController) Review(ctx *gin.Context) {
	currentUser, ok := middleware.CurrentUser(ctx)
	if !ok {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	examID, err := strconv.ParseInt(ctx.Param("examId"), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return
	}
	var request ReviewExamRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.review.Execute(ctx.Request.Context(), currentUser, examID, examusecase.ReviewInput{Status: request.Status, ReviewNote: request.ReviewNote})
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponse(response))
}
func (controller *ExamController) listApprovedResponse(ctx *gin.Context, subjectID *int64) {
	response, err := controller.listApproved.Execute(ctx.Request.Context(), subjectID, ctx.Query("period"))
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponses(response))
}
func parsePathInt64(ctx *gin.Context, parameterName string) (int64, bool) {
	value, err := strconv.ParseInt(ctx.Param(parameterName), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return 0, false
	}
	return value, true
}
func parseRequiredQueryInt64(ctx *gin.Context, parameterName string) (int64, bool) {
	value, err := strconv.ParseInt(ctx.Query(parameterName), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, parameterName+" invalido."))
		return 0, false
	}
	return value, true
}
