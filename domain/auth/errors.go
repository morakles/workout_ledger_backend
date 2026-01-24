package auth

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrStateMismatch = errors.New("oauth state mismatch")
)
