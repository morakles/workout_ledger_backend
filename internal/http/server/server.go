package server

import (
	"net/http"
	token "workout_ledger/internal/auth"
	"workout_ledger/internal/http/handler"
	authmiddleware "workout_ledger/internal/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(paramExerciseHandler *handler.ParamExerciseHandler, paramWorkoutHandler *handler.ParamWorkoutHandler, authHandler *handler.AuthHandler, tokenManager *token.Manager) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger UI
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Get("/auth/google/login", authHandler.StartGoogleLogin)
		r.Get("/auth/google/callback", authHandler.HandleGoogleCallback)
		r.Post("/auth/refresh", authHandler.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.JWTAuth(tokenManager))
			r.Post("/param_exercises", paramExerciseHandler.CreateParamExercise)
			r.Get("/param_exercises", paramExerciseHandler.GetParamExercises)
			r.Get("/param_exercises/{id}", paramExerciseHandler.GetParamExerciseByID)
			r.Post("/param_workouts", paramWorkoutHandler.CreateParamWorkout)
			r.Get("/param_workouts", paramWorkoutHandler.GetParamWorkouts)
			r.Get("/param_workouts/{id}", paramWorkoutHandler.GetParamWorkoutByID)
		})
	})
	return r
}
