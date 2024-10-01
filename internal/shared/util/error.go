package util

import (
	"errors"
	"net/http"
)

// ValidationError Custom error type for validation errors
type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

// NewValidationError Factory function for creating ValidationError
func NewValidationError(err error) error {
	return &ValidationError{Err: err}
}

func HandleError(w http.ResponseWriter, err error, statusCode int) error {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return WriteJson(w, statusCode, ApiError{Error: validationErr.Error()})
	}

	return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	//return WriteJson(w, http.StatusInternalServerError, ApiError{Error: domain.ErrInternalServer.Error()})
}
