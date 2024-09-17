package domain

import "errors"

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrInvalidName      = errors.New("invalid location name")
)
