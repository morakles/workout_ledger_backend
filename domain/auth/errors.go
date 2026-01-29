package auth

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrStateMismatch   = errors.New("oauth state mismatch")
	ErrMissingToken    = errors.New("missing token")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrRefreshInvalid  = errors.New("refresh token invalid")
	ErrRefreshExpired  = errors.New("refresh token expired")
)
