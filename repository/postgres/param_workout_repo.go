package postgres

import (
	"context"
	"database/sql"
	param_workout "workout_ledger/domain/param_workout"

	sq "github.com/Masterminds/squirrel"
)

const workoutTableName = "param_workout"

type ParamWorkoutRepo struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewParamWorkoutRepo(db *sql.DB) *ParamWorkoutRepo {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &ParamWorkoutRepo{
		db: db,
		sb: sb,
	}
}

func (workoutRepo *ParamWorkoutRepo) CreateParamWorkout(ctx context.Context, workout param_workout.ParamWorkout) (param_workout.ParamWorkout, error) {
	sqlInsert := workoutRepo.sb.Insert(workoutTableName).
		Columns("name", "number_of_sets", "rest_between_sets").
		Values(workout.Name, workout.NumberOfSets, sq.Expr("make_interval(secs => ?)", workout.RestBetweenSetsSeconds)).
		Suffix("RETURNING id, name, number_of_sets, EXTRACT(EPOCH FROM rest_between_sets)::int AS rest_between_sets_seconds")

	sqlStr, args, err := sqlInsert.ToSql()
	if err != nil {
		return param_workout.ParamWorkout{}, err
	}

	var created param_workout.ParamWorkout
	err = workoutRepo.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&created.ID,
		&created.Name,
		&created.NumberOfSets,
		&created.RestBetweenSetsSeconds,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return param_workout.ParamWorkout{}, param_workout.ErrConflict
		}
		return param_workout.ParamWorkout{}, err
	}

	return created, nil
}

func (workoutRepo *ParamWorkoutRepo) GetParamWorkoutByID(ctx context.Context, id int64) (param_workout.ParamWorkout, error) {
	sqlSelect := workoutRepo.sb.Select(
		"id",
		"name",
		"number_of_sets",
		"EXTRACT(EPOCH FROM rest_between_sets)::int AS rest_between_sets_seconds",
	).
		From(workoutTableName).
		Where(sq.Eq{"id": id})

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return param_workout.ParamWorkout{}, err
	}

	var workout param_workout.ParamWorkout
	err = workoutRepo.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&workout.ID,
		&workout.Name,
		&workout.NumberOfSets,
		&workout.RestBetweenSetsSeconds,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return param_workout.ParamWorkout{}, param_workout.ErrNotFound
		}
		return param_workout.ParamWorkout{}, err
	}

	return workout, nil
}

func (workoutRepo *ParamWorkoutRepo) GetParamWorkouts(ctx context.Context) ([]param_workout.ParamWorkout, error) {
	sqlSelect := workoutRepo.sb.Select(
		"id",
		"name",
		"number_of_sets",
		"EXTRACT(EPOCH FROM rest_between_sets)::int AS rest_between_sets_seconds",
	).
		From(workoutTableName).
		OrderBy("id")

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := workoutRepo.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []param_workout.ParamWorkout
	for rows.Next() {
		var workout param_workout.ParamWorkout
		if err := rows.Scan(
			&workout.ID,
			&workout.Name,
			&workout.NumberOfSets,
			&workout.RestBetweenSetsSeconds,
		); err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}
