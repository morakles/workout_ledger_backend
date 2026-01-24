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

func isUniqueViolation(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique constraint")
}
