package errs

import "net/http"

type AppError struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e AppError) ShowError() *AppError {
	return &AppError{
		Status:  e.Status,
		Message: e.Message,
	}
}

func NotFoundError(message string, status string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Status:  status,
		Message: message,
	}
}

func UnExpectedError(message string, status string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Status:  status,
		Message: message,
	}
}

func ValidationError(message string, status string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Status:  status,
		Message: message,
	}
}

func AuthenticationError(message string, status string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Status:  status,
		Message: message,
	}
}
