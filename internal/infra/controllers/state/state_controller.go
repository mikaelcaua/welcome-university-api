package statecontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/state"
)

type CreateStateRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type StateController struct {
	list      *stateusecase.ListStatesUseCase
	getByCode *stateusecase.GetStateByCodeUseCase
	create    *stateusecase.CreateStateUseCase
}

func NewStateController(list *stateusecase.ListStatesUseCase, getByCode *stateusecase.GetStateByCodeUseCase, create *stateusecase.CreateStateUseCase) *StateController {
	return &StateController{list: list, getByCode: getByCode, create: create}
}
func (controller *StateController) List(ctx *gin.Context) {
	response, err := controller.list.Execute(ctx.Request.Context())
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *StateController) GetByCode(ctx *gin.Context) {
	response, err := controller.getByCode.Execute(ctx.Request.Context(), ctx.Param("code"))
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *StateController) Create(ctx *gin.Context) {
	var request CreateStateRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.create.Execute(ctx.Request.Context(), request.Code, request.Name)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusCreated, response)
}
