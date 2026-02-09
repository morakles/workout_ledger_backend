package postgres

import (
	"context"
	"database/sql"
	"strings"
	domain "workout_ledger/domain/param_exercise"

	sq "github.com/Masterminds/squirrel"
)

// statics
const exerciseTableName = "param_exercise"

type ParamExerciseRepo struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewParamExerciseRepo(db *sql.DB) *ParamExerciseRepo {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &ParamExerciseRepo{
		db: db,
		sb: sb}
}

func (exerciseRepo *ParamExerciseRepo) CreateParamxercise(ctx context.Context, param_exercise domain.ParamExercise) (int64, error) {
	sqlInsert := exerciseRepo.sb.Insert(exerciseTableName).Columns("exercise_name", "icon_url").
		Values(param_exercise.Name, param_exercise.IconUrl).Suffix("RETURNING id")

	sqlStr, args, err := sqlInsert.ToSql()
	if err != nil {
		return -1, err
	}

	var id int64
	err = exerciseRepo.db.QueryRowContext(ctx, sqlStr, args...).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return -1, domain.ErrExerciseAlreadyExists
		}
		return -1, err
	}
	return id, nil
}

func (exerciseRepo *ParamExerciseRepo) GetParamExercises(ctx context.Context) ([]domain.ParamExercise, error) {
	sqlSelect := exerciseRepo.sb.Select("id", "exercise_name", "icon_url").
		From(exerciseTableName).
		OrderBy("id")

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := exerciseRepo.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []domain.ParamExercise
	for rows.Next() {
		var exercise domain.ParamExercise
		if err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.IconUrl); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exercises, nil
}

func (exerciseRepo *ParamExerciseRepo) GetParamExerciseByID(ctx context.Context, id int64) (domain.ParamExercise, error) {
	sqlSelect := exerciseRepo.sb.Select("id", "exercise_name", "icon_url").
		From(exerciseTableName).
		Where(sq.Eq{"id": id})

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return domain.ParamExercise{}, err
	}

	var exercise domain.ParamExercise
	err = exerciseRepo.db.QueryRowContext(ctx, sqlStr, args...).Scan(&exercise.ID, &exercise.Name, &exercise.IconUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ParamExercise{}, domain.ErrNotFound
		}
		return domain.ParamExercise{}, err
	}

	return exercise, nil
}

func (exerciseRepo *ParamExerciseRepo) UpdateParamExercise(ctx context.Context, param_exercise domain.ParamExercise) (domain.ParamExercise, error) {
	sqlUpdate := exerciseRepo.sb.Update(exerciseTableName).
		Set("exercise_name", param_exercise.Name).
		Set("icon_url", param_exercise.IconUrl).
		Where(sq.Eq{"id": param_exercise.ID}).
		Suffix("RETURNING id, exercise_name, icon_url")

	sqlStr, args, err := sqlUpdate.ToSql()
	if err != nil {
		return domain.ParamExercise{}, err
	}

	var exercise domain.ParamExercise
	err = exerciseRepo.db.QueryRowContext(ctx, sqlStr, args...).Scan(&exercise.ID, &exercise.Name, &exercise.IconUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ParamExercise{}, domain.ErrNotFound
		}
		if isUniqueViolation(err) {
			return domain.ParamExercise{}, domain.ErrExerciseAlreadyExists
		}
		return domain.ParamExercise{}, err
	}

	return exercise, nil
}

func (exerciseRepo *ParamExerciseRepo) DeleteParamExercise(ctx context.Context, id int64) error {
	sqlDelete := exerciseRepo.sb.Delete(exerciseTableName).
		Where(sq.Eq{"id": id})

	sqlStr, args, err := sqlDelete.ToSql()
	if err != nil {
		return err
	}

	result, err := exerciseRepo.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique constraint")
}
