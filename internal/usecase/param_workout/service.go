package param_workout

import (
	"context"
	"strings"
	"workout_ledger/domain/param_workout"
)

type ParamWorkoutRepository interface {
	CreateParamWorkout(ctx context.Context, workout param_workout.ParamWorkout) (param_workout.ParamWorkout, error)
	GetParamWorkoutByID(ctx context.Context, userID, id int64) (param_workout.ParamWorkout, error)
	GetParamWorkouts(ctx context.Context, userID int64) ([]param_workout.ParamWorkout, error)
}

type ParamWorkoutService struct {
	paramWorkoutRepository ParamWorkoutRepository
}

func NewParamWorkoutService(paramWorkoutRepository ParamWorkoutRepository) *ParamWorkoutService {
	return &ParamWorkoutService{paramWorkoutRepository: paramWorkoutRepository}
}

func (s *ParamWorkoutService) CreateParamWorkout(ctx context.Context, dto CreateParamWorkoutDTO) (ParamWorkoutDTO, error) {
	if err := validateParamWorkout(dto); err != nil {
		return ParamWorkoutDTO{}, err
	}

	workout := param_workout.ParamWorkout{
		UserID:                 dto.UserID,
		Name:                   strings.TrimSpace(dto.Name),
		NumberOfSets:           dto.NumberOfSets,
		RestBetweenSetsSeconds: dto.RestBetweenSetsSeconds,
	}

	created, err := s.paramWorkoutRepository.CreateParamWorkout(ctx, workout)
	if err != nil {
		return ParamWorkoutDTO{}, err
	}

	return ParamWorkoutDTO{
		ID:                     created.ID,
		UserID:                 created.UserID,
		Name:                   created.Name,
		NumberOfSets:           created.NumberOfSets,
		RestBetweenSetsSeconds: created.RestBetweenSetsSeconds,
	}, nil
}

func (s *ParamWorkoutService) GetParamWorkoutByID(ctx context.Context, userID, id int64) (ParamWorkoutDTO, error) {
	if userID <= 0 || id <= 0 {
		return ParamWorkoutDTO{}, param_workout.ErrInvalidInput
	}

	workout, err := s.paramWorkoutRepository.GetParamWorkoutByID(ctx, userID, id)
	if err != nil {
		return ParamWorkoutDTO{}, err
	}

	return ParamWorkoutDTO{
		ID:                     workout.ID,
		UserID:                 workout.UserID,
		Name:                   workout.Name,
		NumberOfSets:           workout.NumberOfSets,
		RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
	}, nil
}

func (s *ParamWorkoutService) GetParamWorkouts(ctx context.Context, userID int64) ([]ParamWorkoutDTO, error) {
	if userID <= 0 {
		return nil, param_workout.ErrInvalidInput
	}
	workouts, err := s.paramWorkoutRepository.GetParamWorkouts(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]ParamWorkoutDTO, 0, len(workouts))
	for _, workout := range workouts {
		result = append(result, ParamWorkoutDTO{
			ID:                     workout.ID,
			UserID:                 workout.UserID,
			Name:                   workout.Name,
			NumberOfSets:           workout.NumberOfSets,
			RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
		})
	}

	return result, nil
}

func validateParamWorkout(dto CreateParamWorkoutDTO) error {
	if dto.UserID <= 0 {
		return param_workout.ErrInvalidInput
	}
	name := strings.TrimSpace(dto.Name)
	if name == "" || len(name) > 255 {
		return param_workout.ErrInvalidInput
	}
	if dto.NumberOfSets < 1 {
		return param_workout.ErrInvalidInput
	}
	if dto.RestBetweenSetsSeconds < 0 {
		return param_workout.ErrInvalidInput
	}
	return nil
}
