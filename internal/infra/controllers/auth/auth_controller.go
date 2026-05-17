package authcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/auth"
)

type AuthController struct {
	register *authusecase.RegisterUseCase
	login    *authusecase.LoginUseCase
	refresh  *authusecase.RefreshUseCase
}

func NewAuthController(register *authusecase.RegisterUseCase, login *authusecase.LoginUseCase, refresh *authusecase.RefreshUseCase) *AuthController {
	return &AuthController{register: register, login: login, refresh: refresh}
}
func (controller *AuthController) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	result, err := controller.register.Execute(ctx.Request.Context(), authusecase.RegisterInput{Name: request.Name, Email: request.Email, Password: request.Password})
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusCreated, toResponse(result))
}
func (controller *AuthController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	result, err := controller.login.Execute(ctx.Request.Context(), authusecase.LoginInput{Email: request.Email, Password: request.Password})
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, toResponse(result))
}
func (controller *AuthController) Refresh(ctx *gin.Context) {
	var request RefreshTokenRequest
	if err := httpx.DecodeJSON(ctx, &request); err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	result, err := controller.refresh.Execute(ctx.Request.Context(), authusecase.RefreshInput{RefreshToken: request.RefreshToken})
	if err != nil {
		httpx.RespondError(ctx, err)
		return
	}
	httpx.RespondJSON(ctx, http.StatusOK, toResponse(result))
}
func toResponse(result authusecase.AuthResult) AuthResponse {
	return AuthResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken, TokenType: result.TokenType, ExpiresInSeconds: result.ExpiresInSeconds, User: usercontroller.ToResponse(result.User)}
}
