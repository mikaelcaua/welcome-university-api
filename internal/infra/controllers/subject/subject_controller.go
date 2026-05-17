package subjectcontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/subject"
)

type CreateSubjectRequest struct{ Name string `json:"name"` }
type SubjectController struct {
	list   *subjectusecase.ListSubjectsUseCase
	create *subjectusecase.CreateSubjectUseCase
}

func NewSubjectController(list *subjectusecase.ListSubjectsUseCase, create *subjectusecase.CreateSubjectUseCase) *SubjectController {
	return &SubjectController{list: list, create: create}
}
func (controller *SubjectController) List(ctx *gin.Context) {
	courseID, ok := parsePathInt64(ctx, "courseId")
	if !ok {
		return
	}
	response, err := controller.list.Execute(ctx.Request.Context(), courseID)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *SubjectController) Create(ctx *gin.Context) {
	courseID, ok := parsePathInt64(ctx, "courseId")
	if !ok {
		return
	}
	var request CreateSubjectRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.create.Execute(ctx.Request.Context(), courseID, request.Name)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusCreated, response)
}
func parsePathInt64(ctx *gin.Context, parameterName string) (int64, bool) {
	value, err := strconv.ParseInt(ctx.Param(parameterName), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return 0, false
	}
	return value, true
}
