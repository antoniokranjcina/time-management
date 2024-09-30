package middleware

import "errors"

var (
	ErrNoValidToken            = errors.New("unauthorized: no valid token")
	ErrInvalidToken            = errors.New("unauthorized: invalid token")
	ErrUnExpectedSigningMethod = errors.New("unauthorized: unexpected signing method")
	ErrTokenExpired            = errors.New("unauthorized: token expired")
)
