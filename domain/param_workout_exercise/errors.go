package param_workout_exercise

import "errors"

var (
	ErrWorkoutNotFound    = errors.New("workout not found")
	ErrExerciseNotFound   = errors.New("exercise not found")
	ErrAssignmentNotFound = errors.New("assignment not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("conflict")
)
