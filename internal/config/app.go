package config

import (
	"net/http"
	token "workout_ledger/internal/auth"
	"workout_ledger/internal/db"
	"workout_ledger/internal/http/handler"
	"workout_ledger/internal/http/server"
	auth_uc "workout_ledger/internal/usecase/auth"
	param_exercise_uc "workout_ledger/internal/usecase/param_exercise"
	param_workout_uc "workout_ledger/internal/usecase/param_workout"
	param_workout_exercise_uc "workout_ledger/internal/usecase/param_workout_exercise"
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

	paramExerciseRepo := postgres.NewParamExerciseRepo(pool)
	paramExerciseSvc := param_exercise_uc.NewParamExerciseService(paramExerciseRepo)
	paramExerciseHandler := handler.NewParamExerciseHandler(paramExerciseSvc)

	paramWorkoutRepo := postgres.NewParamWorkoutRepo(pool)
	paramWorkoutSvc := param_workout_uc.NewParamWorkoutService(paramWorkoutRepo)
	paramWorkoutHandler := handler.NewParamWorkoutHandler(paramWorkoutSvc)

	paramWorkoutExerciseRepo := postgres.NewParamWorkoutExerciseRepo(pool)
	paramWorkoutExerciseSvc := param_workout_exercise_uc.NewParamWorkoutExerciseService(paramWorkoutExerciseRepo)
	paramWorkoutExerciseHandler := handler.NewParamWorkoutExerciseHandler(paramWorkoutExerciseSvc)

	tokenManager, err := token.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	if err != nil {
		return nil, err
	}

	authRepo := postgres.NewAuthRepo(pool)
	authSvc := auth_uc.NewAuthService(authRepo, tokenManager)
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

	httpHandler := server.New(paramExerciseHandler, paramWorkoutHandler, paramWorkoutExerciseHandler, authHandler, tokenManager)

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
