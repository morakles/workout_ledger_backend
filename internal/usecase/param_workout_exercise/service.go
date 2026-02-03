package param_workout_exercise

import (
	"context"
	"workout_ledger/domain/param_workout_exercise"
)

type ParamWorkoutExerciseRepository interface {
	AddExerciseToWorkout(ctx context.Context, workoutID, exerciseID int64, order int) (param_workout_exercise.WorkoutExercise, error)
	ListWorkoutExercises(ctx context.Context, workoutID int64) ([]param_workout_exercise.WorkoutExercise, error)
	RemoveExerciseFromWorkout(ctx context.Context, workoutID, exerciseID int64) error
}

type ParamWorkoutExerciseService struct {
	paramWorkoutExerciseRepository ParamWorkoutExerciseRepository
}

func NewParamWorkoutExerciseService(paramWorkoutExerciseRepository ParamWorkoutExerciseRepository) *ParamWorkoutExerciseService {
	return &ParamWorkoutExerciseService{paramWorkoutExerciseRepository: paramWorkoutExerciseRepository}
}

func (s *ParamWorkoutExerciseService) AddExerciseToWorkout(ctx context.Context, dto AddWorkoutExerciseDTO) (WorkoutExerciseDTO, error) {
	if err := validateWorkoutExerciseInput(dto.WorkoutID, dto.ExerciseID, dto.ExerciseOrder); err != nil {
		return WorkoutExerciseDTO{}, err
	}

	assignment, err := s.paramWorkoutExerciseRepository.AddExerciseToWorkout(ctx, dto.WorkoutID, dto.ExerciseID, dto.ExerciseOrder)
	if err != nil {
		return WorkoutExerciseDTO{}, err
	}

	return WorkoutExerciseDTO{
		WorkoutID:     assignment.WorkoutID,
		ExerciseID:    assignment.ExerciseID,
		ExerciseOrder: assignment.ExerciseOrder,
		ExerciseName:  assignment.ExerciseName,
	}, nil
}

func (s *ParamWorkoutExerciseService) ListWorkoutExercises(ctx context.Context, workoutID int64) ([]WorkoutExerciseDTO, error) {
	if workoutID <= 0 {
		return nil, param_workout_exercise.ErrInvalidInput
	}

	assignments, err := s.paramWorkoutExerciseRepository.ListWorkoutExercises(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	result := make([]WorkoutExerciseDTO, 0, len(assignments))
	for _, assignment := range assignments {
		result = append(result, WorkoutExerciseDTO{
			WorkoutID:     assignment.WorkoutID,
			ExerciseID:    assignment.ExerciseID,
			ExerciseOrder: assignment.ExerciseOrder,
			ExerciseName:  assignment.ExerciseName,
		})
	}

	return result, nil
}

func (s *ParamWorkoutExerciseService) RemoveExerciseFromWorkout(ctx context.Context, workoutID, exerciseID int64) error {
	if workoutID <= 0 || exerciseID <= 0 {
		return param_workout_exercise.ErrInvalidInput
	}

	return s.paramWorkoutExerciseRepository.RemoveExerciseFromWorkout(ctx, workoutID, exerciseID)
}

func validateWorkoutExerciseInput(workoutID, exerciseID int64, order int) error {
	if workoutID <= 0 || exerciseID <= 0 || order < 1 {
		return param_workout_exercise.ErrInvalidInput
	}
	return nil
}
