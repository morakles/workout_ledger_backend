package param_workout

import "errors"

var (
	ErrNotFound     = errors.New("workout not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)
