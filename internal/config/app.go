package config

import (
	"net/http"
	"workout_ledger/internal/db"
	"workout_ledger/internal/http/handler"
	"workout_ledger/internal/http/server"
	param_exercise_uc "workout_ledger/internal/usecase/param_exercise"
	"workout_ledger/repository/postgres"
)

type App struct {
	cfg    Config
	server *http.Server
}

func New(cfg Config) (*App, error) {
	pool, err := db.New(cfg.DSN)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewParamExerciseRepo(pool)
	svc := param_exercise_uc.NewParamExerciseService(repo)
	h := handler.NewParamExerciseHandler(svc)

	httpHandler := server.New(h)

	srv := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: httpHandler,
	}
	return &App{
		cfg:    cfg,
		server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}
