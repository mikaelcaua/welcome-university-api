package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (httpError HTTPError) Error() string {
	return httpError.Message
}

func NewHTTPError(statusCode int, message string) HTTPError {
	return HTTPError{StatusCode: statusCode, Message: message}
}

func RespondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	var httpError HTTPError
	if errors.As(err, &httpError) {
		RespondJSON(w, httpError.StatusCode, map[string]string{"message": httpError.Message})
		return
	}
	RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Erro interno."})
}

func DecodeJSON(r *http.Request, destination any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(destination); err != nil {
		return NewHTTPError(http.StatusBadRequest, "JSON invalido.")
	}
	return nil
}
