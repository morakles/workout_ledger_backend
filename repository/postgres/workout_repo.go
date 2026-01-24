package postgres

import "database/sql"

type ParamExerciseRepository struct {
	db *sql.DB
}

func NewParamExerciseRepository(db *sql.DB) *ParamExerciseRepository {
	return &ParamExerciseRepository{db: db}
}
