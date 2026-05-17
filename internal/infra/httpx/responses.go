package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (httpError HTTPError) Error() string { return httpError.Message }
func NewHTTPError(statusCode int, message string) HTTPError {
	return HTTPError{StatusCode: statusCode, Message: message}
}
func RespondJSON(ctx *gin.Context, statusCode int, payload any) { ctx.JSON(statusCode, payload) }
func RespondError(ctx *gin.Context, err error) {
	var httpError HTTPError
	if errors.As(err, &httpError) {
		RespondJSON(ctx, httpError.StatusCode, gin.H{"message": httpError.Message})
		return
	}
	var applicationError domainerrors.ApplicationError
	if errors.As(err, &applicationError) {
		statusCode := http.StatusInternalServerError
		switch applicationError.Kind {
		case domainerrors.KindValidation:
			statusCode = http.StatusBadRequest
		case domainerrors.KindUnauthorized:
			statusCode = http.StatusUnauthorized
		case domainerrors.KindNotFound:
			statusCode = http.StatusNotFound
		case domainerrors.KindConflict:
			statusCode = http.StatusConflict
		}
		RespondJSON(ctx, statusCode, gin.H{"message": applicationError.Message})
		return
	}
	RespondJSON(ctx, http.StatusInternalServerError, gin.H{"message": "Erro interno."})
}
func DecodeJSON(ctx *gin.Context, destination any) error {
	if err := ctx.ShouldBindJSON(destination); err != nil {
		return NewHTTPError(http.StatusBadRequest, "JSON invalido.")
	}
	return nil
}
