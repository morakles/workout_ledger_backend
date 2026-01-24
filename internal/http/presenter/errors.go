package presenter

import (
	"errors"
	"net/http"
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
	default:
		return http.StatusInternalServerError, APIError{Code: "INTERNAL", Message: "Internal server error"}
	}
}
