package coursecontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/course"
)

type CreateCourseRequest struct{ Name string `json:"name"` }
type CourseController struct {
	list   *courseusecase.ListCoursesUseCase
	create *courseusecase.CreateCourseUseCase
}

func NewCourseController(list *courseusecase.ListCoursesUseCase, create *courseusecase.CreateCourseUseCase) *CourseController {
	return &CourseController{list: list, create: create}
}
func (controller *CourseController) List(ctx *gin.Context) {
	universityID, ok := parsePathInt64(ctx, "universityId")
	if !ok {
		return
	}
	response, err := controller.list.Execute(ctx.Request.Context(), universityID)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *CourseController) Create(ctx *gin.Context) {
	universityID, ok := parsePathInt64(ctx, "universityId")
	if !ok {
		return
	}
	var request CreateCourseRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.create.Execute(ctx.Request.Context(), universityID, request.Name)
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
