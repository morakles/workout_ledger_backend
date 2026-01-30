package param_workout

import (
	"context"
	"strings"
	"workout_ledger/domain/param_workout"
)

type ParamWorkoutRepository interface {
	CreateParamWorkout(ctx context.Context, workout param_workout.ParamWorkout) (param_workout.ParamWorkout, error)
	GetParamWorkoutByID(ctx context.Context, id int64) (param_workout.ParamWorkout, error)
	GetParamWorkouts(ctx context.Context) ([]param_workout.ParamWorkout, error)
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
		Name:                   created.Name,
		NumberOfSets:           created.NumberOfSets,
		RestBetweenSetsSeconds: created.RestBetweenSetsSeconds,
	}, nil
}

func (s *ParamWorkoutService) GetParamWorkoutByID(ctx context.Context, id int64) (ParamWorkoutDTO, error) {
	if id <= 0 {
		return ParamWorkoutDTO{}, param_workout.ErrInvalidInput
	}

	workout, err := s.paramWorkoutRepository.GetParamWorkoutByID(ctx, id)
	if err != nil {
		return ParamWorkoutDTO{}, err
	}

	return ParamWorkoutDTO{
		ID:                     workout.ID,
		Name:                   workout.Name,
		NumberOfSets:           workout.NumberOfSets,
		RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
	}, nil
}

func (s *ParamWorkoutService) GetParamWorkouts(ctx context.Context) ([]ParamWorkoutDTO, error) {
	workouts, err := s.paramWorkoutRepository.GetParamWorkouts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ParamWorkoutDTO, 0, len(workouts))
	for _, workout := range workouts {
		result = append(result, ParamWorkoutDTO{
			ID:                     workout.ID,
			Name:                   workout.Name,
			NumberOfSets:           workout.NumberOfSets,
			RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
		})
	}

	return result, nil
}

func validateParamWorkout(dto CreateParamWorkoutDTO) error {
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
