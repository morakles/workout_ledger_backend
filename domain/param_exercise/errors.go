package param_exercise

import "errors"

var (
	ErrExerciseAlreadyExists = errors.New("ExerciseAlreadyExists")
	ErrNotFound              = errors.New("Exercise not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrForbidden             = errors.New("forbidden")
	ErrWorkoutNotActive      = errors.New("workout not active")
	ErrConflict              = errors.New("conflict")
)
