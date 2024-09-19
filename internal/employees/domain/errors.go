package domain

import "errors"

var (
	ErrFirstNameTooShort = errors.New("first name is too short")
	ErrLastNameTooShort  = errors.New("last name is too short")
	ErrEmailTaken        = errors.New("email already taken")
	ErrEmailWrongFormat  = errors.New("email has wrong format")
	ErrPasswordTooShort  = errors.New("password is too short")
	ErrEmployeeNotFound  = errors.New("employee not found")
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
