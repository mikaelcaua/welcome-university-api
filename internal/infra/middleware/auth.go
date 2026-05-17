package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
)

const authenticatedUserContextKey = "authenticatedUser"

type AccessTokenValidator interface {
	ExtractEmailFromAccessToken(tokenText string) (string, error)
}
type UserFinder interface {
	FindByEmail(ctx context.Context, email string) (user.User, bool, error)
}

func Authentication(tokenValidator AccessTokenValidator, userFinder UserFinder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerValue := ctx.GetHeader("Authorization")
		tokenText := strings.TrimPrefix(headerValue, "Bearer ")
		if tokenText == "" || tokenText == headerValue {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
			ctx.Abort()
			return
		}
		email, err := tokenValidator.ExtractEmailFromAccessToken(tokenText)
		if err != nil {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Token invalido."))
			ctx.Abort()
			return
		}
		user, found, err := userFinder.FindByEmail(ctx.Request.Context(), email)
		if err != nil {
			httpx.RespondError(ctx, err)
			ctx.Abort()
			return
		}
		if !found {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Usuario nao encontrado."))
			ctx.Abort()
			return
		}
		ctx.Set(authenticatedUserContextKey, user)
		ctx.Next()
	}
}

func RequireRoles(roles ...user.Role) gin.HandlerFunc {
	allowedRoles := map[user.Role]bool{}
	for _, role := range roles {
		allowedRoles[role] = true
	}
	return func(ctx *gin.Context) {
		user, ok := CurrentUser(ctx)
		if !ok {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
			ctx.Abort()
			return
		}
		if !allowedRoles[user.Role] {
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusForbidden, "Acesso negado."))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CurrentUser(ctx *gin.Context) (user.User, bool) {
	value, ok := ctx.Get(authenticatedUserContextKey)
	if !ok {
		return user.User{}, false
	}
	user, ok := value.(user.User)
	return user, ok
}
