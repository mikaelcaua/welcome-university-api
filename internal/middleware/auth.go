package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
)

type contextKey string

const authenticatedUserContextKey contextKey = "authenticatedUser"

type AccessTokenValidator interface {
	ExtractEmailFromAccessToken(tokenText string) (string, error)
}

type UserFinder interface {
	FindByEmail(ctx context.Context, email string) (models.User, bool, error)
}

func Authentication(tokenValidator AccessTokenValidator, userFinder UserFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerValue := r.Header.Get("Authorization")
			tokenText := strings.TrimPrefix(headerValue, "Bearer ")
			if tokenText == "" || tokenText == headerValue {
				httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
				return
			}

			email, err := tokenValidator.ExtractEmailFromAccessToken(tokenText)
			if err != nil {
				httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Token invalido."))
				return
			}
			user, found, err := userFinder.FindByEmail(r.Context(), email)
			if err != nil {
				httpx.RespondError(w, err)
				return
			}
			if !found {
				httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Usuario nao encontrado."))
				return
			}
			ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...models.Role) func(http.Handler) http.Handler {
	allowedRoles := map[models.Role]bool{}
	for _, role := range roles {
		allowedRoles[role] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := CurrentUser(r.Context())
			if !ok {
				httpx.RespondError(w, httpx.NewHTTPError(http.StatusUnauthorized, "Autenticacao obrigatoria."))
				return
			}
			if !allowedRoles[user.Role] {
				httpx.RespondError(w, httpx.NewHTTPError(http.StatusForbidden, "Acesso negado."))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CurrentUser(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(authenticatedUserContextKey).(models.User)
	return user, ok
}
