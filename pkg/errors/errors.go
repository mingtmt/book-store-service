package errors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Predefined errors
var (
	ErrNotFound          = New("NOT_FOUND", "resource not found", http.StatusNotFound)
	ErrBadRequest        = New("BAD_REQUEST", "invalid request", http.StatusBadRequest)
	ErrUnauthorized      = New("UNAUTHORIZED", "unauthorized access", http.StatusUnauthorized)
	ErrInternal          = New("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
	ErrUserNotFound      = New("USER_NOT_FOUND", "user not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("USER_ALREADY_EXISTS", "user already exists", http.StatusConflict)
	ErrInvalidPassword   = New("INVALID_PASSWORD", "invalid password", http.StatusUnauthorized)
	ErrInvalidToken      = New("INVALID_TOKEN", "invalid token", http.StatusUnauthorized)
)
