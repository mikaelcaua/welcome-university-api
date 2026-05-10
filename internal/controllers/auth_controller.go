package controllers

import (
	"net/http"

	"github.com/mikaelcaua/welcome-university-api/internal/dto"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	response, err := controller.authService.Register(r.Context(), request)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusCreated, response)
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	response, err := controller.authService.Login(r.Context(), request)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}

func (controller *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var request dto.RefreshTokenRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.RespondError(w, err)
		return
	}
	response, err := controller.authService.Refresh(r.Context(), request)
	if err != nil {
		httpx.RespondError(w, err)
		return
	}
	httpx.RespondJSON(w, http.StatusOK, response)
}
