package postgres

import (
	"context"
	"database/sql"
	"workout_ledger/domain/param_workout_exercise"

	sq "github.com/Masterminds/squirrel"
)

const (
	workoutExercisesTableName = "exercises_in_workout"
	paramExerciseTableName    = "param_exercise"
	paramWorkoutTableName     = "param_workout"
)

type ParamWorkoutExerciseRepo struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewParamWorkoutExerciseRepo(db *sql.DB) *ParamWorkoutExerciseRepo {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &ParamWorkoutExerciseRepo{db: db, sb: sb}
}

func (repo *ParamWorkoutExerciseRepo) AddExerciseToWorkout(ctx context.Context, workoutID, exerciseID int64, order int, defaultSets, defaultRestSeconds, defaultReps *int) (param_workout_exercise.WorkoutExercise, error) {
	workoutExists, err := repo.workoutExists(ctx, workoutID)
	if err != nil {
		return param_workout_exercise.WorkoutExercise{}, err
	}
	if !workoutExists {
		return param_workout_exercise.WorkoutExercise{}, param_workout_exercise.ErrWorkoutNotFound
	}

	exerciseExists, err := repo.exerciseExists(ctx, exerciseID)
	if err != nil {
		return param_workout_exercise.WorkoutExercise{}, err
	}
	if !exerciseExists {
		return param_workout_exercise.WorkoutExercise{}, param_workout_exercise.ErrExerciseNotFound
	}

	defaultSetsValue := toNullInt16(defaultSets)
	defaultRepsValue := toNullInt16(defaultReps)

	sqlInsert := repo.sb.Insert(workoutExercisesTableName).
		Columns("param_workout_id", "param_exercise_id", "default_sets", "default_rest", "default_reps", "order_index").
		Values(workoutID, exerciseID, defaultSetsValue, toIntervalExpr(defaultRestSeconds), defaultRepsValue, order).
		Suffix("RETURNING param_workout_id, param_exercise_id, order_index, (SELECT exercise_name FROM param_exercise WHERE id = param_exercise_id) AS exercise_name, default_sets, EXTRACT(EPOCH FROM default_rest)::int AS default_rest_seconds, default_reps")

	sqlStr, args, err := sqlInsert.ToSql()
	if err != nil {
		return param_workout_exercise.WorkoutExercise{}, err
	}

	var assignment param_workout_exercise.WorkoutExercise
	var defaultSetsResult sql.NullInt16
	var defaultRestResult sql.NullInt64
	var defaultRepsResult sql.NullInt16
	err = repo.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&assignment.WorkoutID,
		&assignment.ExerciseID,
		&assignment.ExerciseOrder,
		&assignment.ExerciseName,
		&defaultSetsResult,
		&defaultRestResult,
		&defaultRepsResult,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return param_workout_exercise.WorkoutExercise{}, param_workout_exercise.ErrConflict
		}
		return param_workout_exercise.WorkoutExercise{}, err
	}

	if defaultSetsResult.Valid {
		v := int(defaultSetsResult.Int16)
		assignment.DefaultSets = &v
	}
	if defaultRestResult.Valid {
		v := int(defaultRestResult.Int64)
		assignment.DefaultRestSeconds = &v
	}
	if defaultRepsResult.Valid {
		v := int(defaultRepsResult.Int16)
		assignment.DefaultReps = &v
	}

	return assignment, nil
}

func (repo *ParamWorkoutExerciseRepo) ListWorkoutExercises(ctx context.Context, workoutID int64) ([]param_workout_exercise.WorkoutExercise, error) {
	workoutExists, err := repo.workoutExists(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if !workoutExists {
		return nil, param_workout_exercise.ErrWorkoutNotFound
	}

	sqlSelect := repo.sb.Select(
		"eiw.param_workout_id",
		"eiw.param_exercise_id",
		"eiw.order_index",
		"pe.exercise_name AS exercise_name",
		"eiw.default_sets",
		"EXTRACT(EPOCH FROM eiw.default_rest)::int AS default_rest_seconds",
		"eiw.default_reps",
	).
		From(workoutExercisesTableName + " eiw").
		Join(paramExerciseTableName + " pe ON pe.id = eiw.param_exercise_id").
		Where(sq.Eq{"eiw.param_workout_id": workoutID}).
		OrderBy("eiw.order_index ASC")

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make([]param_workout_exercise.WorkoutExercise, 0)
	for rows.Next() {
		var assignment param_workout_exercise.WorkoutExercise
		var defaultSetsResult sql.NullInt16
		var defaultRestResult sql.NullInt64
		var defaultRepsResult sql.NullInt16
		if err := rows.Scan(
			&assignment.WorkoutID,
			&assignment.ExerciseID,
			&assignment.ExerciseOrder,
			&assignment.ExerciseName,
			&defaultSetsResult,
			&defaultRestResult,
			&defaultRepsResult,
		); err != nil {
			return nil, err
		}
		if defaultSetsResult.Valid {
			v := int(defaultSetsResult.Int16)
			assignment.DefaultSets = &v
		}
		if defaultRestResult.Valid {
			v := int(defaultRestResult.Int64)
			assignment.DefaultRestSeconds = &v
		}
		if defaultRepsResult.Valid {
			v := int(defaultRepsResult.Int16)
			assignment.DefaultReps = &v
		}
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (repo *ParamWorkoutExerciseRepo) RemoveExerciseFromWorkout(ctx context.Context, workoutID, exerciseID int64) error {
	workoutExists, err := repo.workoutExists(ctx, workoutID)
	if err != nil {
		return err
	}
	if !workoutExists {
		return param_workout_exercise.ErrWorkoutNotFound
	}

	sqlDelete := repo.sb.Delete(workoutExercisesTableName).
		Where(sq.Eq{"param_workout_id": workoutID, "param_exercise_id": exerciseID})

	sqlStr, args, err := sqlDelete.ToSql()
	if err != nil {
		return err
	}

	result, err := repo.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return param_workout_exercise.ErrAssignmentNotFound
	}

	return nil
}

func (repo *ParamWorkoutExerciseRepo) workoutExists(ctx context.Context, workoutID int64) (bool, error) {
	sqlSelect := repo.sb.Select("1").
		From(paramWorkoutTableName).
		Where(sq.Eq{"id": workoutID}).
		Limit(1)

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = repo.db.QueryRowContext(ctx, sqlStr, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *ParamWorkoutExerciseRepo) exerciseExists(ctx context.Context, exerciseID int64) (bool, error) {
	sqlSelect := repo.sb.Select("1").
		From(paramExerciseTableName).
		Where(sq.Eq{"id": exerciseID}).
		Limit(1)

	sqlStr, args, err := sqlSelect.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = repo.db.QueryRowContext(ctx, sqlStr, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func toNullInt16(value *int) sql.NullInt16 {
	if value == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: int16(*value), Valid: true}
}

func toIntervalExpr(value *int) interface{} {
	if value == nil {
		return nil
	}
	return sq.Expr("make_interval(secs => ?)", *value)
}
