package usercontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/user"
)

type UserController struct {
	getCurrentUser *userusecase.GetCurrentUserUseCase
	listUsers      *userusecase.ListUsersUseCase
	updateRole     *userusecase.UpdateUserRoleUseCase
}

func NewUserController(getCurrentUser *userusecase.GetCurrentUserUseCase, listUsers *userusecase.ListUsersUseCase, updateRole *userusecase.UpdateUserRoleUseCase) *UserController {
	return &UserController{getCurrentUser: getCurrentUser, listUsers: listUsers, updateRole: updateRole}
}
func (controller *UserController) Me(ctx *gin.Context) {
	currentUser, ok := middleware.CurrentUser(ctx)
	if !ok {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponse(controller.getCurrentUser.Execute(currentUser)))
}
func (controller *UserController) ListAll(ctx *gin.Context) {
	users, err := controller.listUsers.Execute(ctx.Request.Context())
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response := make([]UserResponse, 0, len(users))
	for _, item := range users {
		response = append(response, ToResponse(item))
	}
	httpx.RespondJSON(ctx, http.StatusOK, response)
}
func (controller *UserController) UpdateRole(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return
	}
	var request UpdateUserRoleRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	response, err := controller.updateRole.Execute(ctx.Request.Context(), userID, request.Role)
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, ToResponse(response))
}
