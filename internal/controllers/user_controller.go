package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) Me(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r.Context())
	if !ok {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
		return
	}
	httpx.RespondJSON(w, http.StatusOK, controller.userService.Me(currentUser))
}

func (controller *UserController) ListAll(w http.ResponseWriter, r *http.Request) {
	users, err := controller.userService.ListAll(r.Context())
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, users)
}

func (controller *UserController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.RespondError(w, httpx.NewHTTPError(http.StatusBadRequest, "Parametro invalido."))
		return
	}
	var request dto.UpdateUserRoleRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	response, err := controller.userService.UpdateRole(r.Context(), userID, request.Role)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}
