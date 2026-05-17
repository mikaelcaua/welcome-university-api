package domainerrors

type Kind string

const (
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
)

type ApplicationError struct {
	Kind    Kind
	Message string
}

func (err ApplicationError) Error() string { return err.Message }

func Validation(message string) ApplicationError {
	return ApplicationError{Kind: KindValidation, Message: message}
}

func Unauthorized(message string) ApplicationError {
	return ApplicationError{Kind: KindUnauthorized, Message: message}
}

func NotFound(message string) ApplicationError {
	return ApplicationError{Kind: KindNotFound, Message: message}
}

func Conflict(message string) ApplicationError {
	return ApplicationError{Kind: KindConflict, Message: message}
}
