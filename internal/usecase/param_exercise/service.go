package param_exercise

import (
	"context"
	domain "workout_ledger/domain/param_exercise"
)

type ParamExerciseRepository interface {
	CreateParamxercise(ctx context.Context, param_exercise domain.ParamExercise) (int64, error)
	GetParamExercises(ctx context.Context) ([]domain.ParamExercise, error)
	GetParamExerciseByID(ctx context.Context, id int64) (domain.ParamExercise, error)
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

func (s *ParamExerciseService) GetParamExercises(ctx context.Context) ([]ParamExerciseDTO, error) {
	paramExercises, err := s.paramExerciseRepository.GetParamExercises(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ParamExerciseDTO, 0, len(paramExercises))
	for _, exercise := range paramExercises {
		result = append(result, ParamExerciseDTO{
			ID:      exercise.ID,
			Name:    exercise.Name,
			IconUrl: exercise.IconUrl,
		})
	}

	return result, nil
}

func (s *ParamExerciseService) GetParamExerciseByID(ctx context.Context, id int64) (ParamExerciseDTO, error) {
	if id <= 0 {
		return ParamExerciseDTO{}, domain.ErrInvalidInput
	}

	exercise, err := s.paramExerciseRepository.GetParamExerciseByID(ctx, id)
	if err != nil {
		return ParamExerciseDTO{}, err
	}

	return ParamExerciseDTO{
		ID:      exercise.ID,
		Name:    exercise.Name,
		IconUrl: exercise.IconUrl,
	}, nil
}
