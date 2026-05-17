package universitycontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/university"
)

type CreateUniversityRequest struct{ Name string `json:"name"` }
type UniversityController struct {
	list   *universityusecase.ListUniversitiesUseCase
	create *universityusecase.CreateUniversityUseCase
}

func NewUniversityController(list *universityusecase.ListUniversitiesUseCase, create *universityusecase.CreateUniversityUseCase) *UniversityController {
	return &UniversityController{list: list, create: create}
}
func (controller *UniversityController) List(ctx *gin.Context) {
	stateID, ok := parsePathInt64(ctx, "stateId")
	if !ok {
		return
	}
	response, err := controller.list.Execute(ctx.Request.Context(), stateID)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *UniversityController) Create(ctx *gin.Context) {
	stateID, ok := parsePathInt64(ctx, "stateId")
	if !ok {
		return
	}
	var request CreateUniversityRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.create.Execute(ctx.Request.Context(), stateID, request.Name)
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
