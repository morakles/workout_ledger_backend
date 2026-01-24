package config

import (
	"net/http"
	"workout_ledger/internal/db"
	"workout_ledger/internal/http/handler"
	"workout_ledger/internal/http/server"
	auth_uc "workout_ledger/internal/usecase/auth"
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

	authRepo := postgres.NewAuthRepo(pool)
	authSvc := auth_uc.NewAuthService(authRepo)
	var googleOAuthConfig *handler.GoogleOAuthConfig
	if cfg.GoogleOAuth != nil {
		googleOAuthConfig = &handler.GoogleOAuthConfig{
			ClientID:      cfg.GoogleOAuth.ClientID,
			ClientSecret:  cfg.GoogleOAuth.ClientSecret,
			RedirectURL:   cfg.GoogleOAuth.RedirectURL,
			AllowedDomain: cfg.GoogleOAuth.AllowedDomain,
		}
	}
	authHandler := handler.NewAuthHandler(authSvc, googleOAuthConfig)

	httpHandler := server.New(h, authHandler)

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
