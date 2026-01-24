package param_exercise

import (
	"context"
	domain "workout_ledger/domain/param_exercise"
)

type ParamExerciseRepository interface {
	CreateParamxercise(ctx context.Context, param_exercise domain.ParamExercise) (int64, error)
}

type ParamExerciseService struct {
	paramExerciseRepository ParamExerciseRepository
}

func NewParamExerciseService(paramExerciseRepository ParamExerciseRepository) *ParamExerciseService {
	return &ParamExerciseService{
		paramExerciseRepository: paramExerciseRepository,
	}
}

func (s *ParamExerciseService) CreateParamExercise(ctx context.Context, dto CreateParamExerciseDTO) (CreateParamExerciseResultDTO, error) {

	paramExercise := domain.ParamExercise{
		Name:    dto.Name,
		IconUrl: dto.IconUrl,
	}

	id, err := s.paramExerciseRepository.CreateParamxercise(ctx, paramExercise)
	if err != nil {
		return CreateParamExerciseResultDTO{}, err
	}

	return CreateParamExerciseResultDTO{ID: id}, nil
}
