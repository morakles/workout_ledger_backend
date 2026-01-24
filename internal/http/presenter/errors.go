package presenter

import (
	"errors"
	"net/http"
	authDomain "workout_ledger/domain/auth"
	paramExerciseDomain "workout_ledger/domain/param_exercise"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func StatusAndError(err error) (int, APIError) {
	switch {
	case errors.Is(err, paramExerciseDomain.ErrInvalidInput):
		return http.StatusBadRequest, APIError{Code: "INVALID_INPUT", Message: "Invalid input"}
	case errors.Is(err, paramExerciseDomain.ErrForbidden):
		return http.StatusForbidden, APIError{Code: "FORBIDDEN", Message: "Forbidden"}
	case errors.Is(err, paramExerciseDomain.ErrNotFound):
		return http.StatusNotFound, APIError{Code: "NOT_FOUND", Message: "Not found"}
	case errors.Is(err, paramExerciseDomain.ErrWorkoutNotActive):
		return http.StatusConflict, APIError{Code: "WORKOUT_NOT_ACTIVE", Message: "Workout is not active"}
	case errors.Is(err, paramExerciseDomain.ErrConflict):
		return http.StatusConflict, APIError{Code: "CONFLICT", Message: "Conflict"}
	case errors.Is(err, authDomain.ErrInvalidInput):
		return http.StatusBadRequest, APIError{Code: "INVALID_INPUT", Message: "Invalid input"}
	case errors.Is(err, authDomain.ErrUnauthorized):
		return http.StatusUnauthorized, APIError{Code: "UNAUTHORIZED", Message: "Unauthorized"}
	case errors.Is(err, authDomain.ErrStateMismatch):
		return http.StatusBadRequest, APIError{Code: "STATE_MISMATCH", Message: "OAuth state mismatch"}
	default:
		return http.StatusInternalServerError, APIError{Code: "INTERNAL", Message: "Internal server error"}
	}
}
