package param_exercise

import (
	"context"
	"testing"
	"workout_ledger/domain/param_exercise"
)

type fakeParamExerciseRepo struct {
	getByID func(ctx context.Context, id int64) (param_exercise.ParamExercise, error)
	update  func(ctx context.Context, exercise param_exercise.ParamExercise) (param_exercise.ParamExercise, error)
	delete  func(ctx context.Context, id int64) error
}

func (f fakeParamExerciseRepo) CreateParamxercise(ctx context.Context, exercise param_exercise.ParamExercise) (int64, error) {
	return 0, nil
}

func (f fakeParamExerciseRepo) GetParamExercises(ctx context.Context) ([]param_exercise.ParamExercise, error) {
	return nil, nil
}

func (f fakeParamExerciseRepo) GetParamExerciseByID(ctx context.Context, id int64) (param_exercise.ParamExercise, error) {
	if f.getByID == nil {
		return param_exercise.ParamExercise{}, nil
	}
	return f.getByID(ctx, id)
}

func (f fakeParamExerciseRepo) UpdateParamExercise(ctx context.Context, exercise param_exercise.ParamExercise) (param_exercise.ParamExercise, error) {
	if f.update == nil {
		return exercise, nil
	}
	return f.update(ctx, exercise)
}

func (f fakeParamExerciseRepo) DeleteParamExercise(ctx context.Context, id int64) error {
	if f.delete == nil {
		return nil
	}
	return f.delete(ctx, id)
}

func TestGetParamExerciseByIDRejectsNonPositiveID(t *testing.T) {
	t.Parallel()

	svc := NewParamExerciseService(fakeParamExerciseRepo{})

	_, err := svc.GetParamExerciseByID(context.Background(), 0)
	if err != param_exercise.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetParamExerciseByIDReturnsDTO(t *testing.T) {
	t.Parallel()

	svc := NewParamExerciseService(fakeParamExerciseRepo{
		getByID: func(ctx context.Context, id int64) (param_exercise.ParamExercise, error) {
			return param_exercise.ParamExercise{
				ID:      id,
				Name:    "Push-up",
				IconUrl: "https://example.com/icon.png",
			}, nil
		},
	})

	dto, err := svc.GetParamExerciseByID(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.ID != 2 || dto.Name != "Push-up" || dto.IconUrl != "https://example.com/icon.png" {
		t.Fatalf("unexpected dto: %#v", dto)
	}
}

func TestGetParamExerciseByIDPropagatesError(t *testing.T) {
	t.Parallel()

	svc := NewParamExerciseService(fakeParamExerciseRepo{
		getByID: func(ctx context.Context, id int64) (param_exercise.ParamExercise, error) {
			return param_exercise.ParamExercise{}, param_exercise.ErrNotFound
		},
	})

	_, err := svc.GetParamExerciseByID(context.Background(), 3)
	if err != param_exercise.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
