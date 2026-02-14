package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workout_ledger/domain/param_workout_exercise"
	"workout_ledger/internal/http/presenter"
	param_workout_exercise_uc "workout_ledger/internal/usecase/param_workout_exercise"

	"github.com/go-chi/chi/v5"
)

type ParamWorkoutExerciseHandler struct {
	svc param_workout_exercise_uc.ParamWorkoutExerciseService
}

func NewParamWorkoutExerciseHandler(svc *param_workout_exercise_uc.ParamWorkoutExerciseService) *ParamWorkoutExerciseHandler {
	return &ParamWorkoutExerciseHandler{svc: *svc}
}

// swagger:parameters addExerciseToWorkout
// swagger:parameters listWorkoutExercises
type workoutIDParam struct {
	// Workout ID
	//
	// in: path
	WorkoutID int64 `json:"workoutId"`
}

type removeExerciseFromWorkoutParams struct {
	// Workout ID
	//
	// in: path
	WorkoutID int64 `json:"workoutId"`
	// Exercise ID
	//
	// in: path
	ExerciseID int64 `json:"exerciseId"`
}

// swagger:parameters addExerciseToWorkout
// in: body
// required: true
// body to add exercise to workout
type addExerciseToWorkoutRequest struct {
	ExerciseID    int64 `json:"exercise_id" example:"123"`
	ExerciseOrder int   `json:"exercise_order" example:"1"`
	DefaultSets   *int  `json:"default_sets,omitempty" example:"4"`
	DefaultRest   *int  `json:"default_rest,omitempty" example:"90"`
	DefaultReps   *int  `json:"default_reps,omitempty" example:"10"`
}

type workoutExerciseResponse struct {
	WorkoutID     int64  `json:"workout_id" example:"1"`
	ExerciseID    int64  `json:"exercise_id" example:"123"`
	ExerciseOrder int    `json:"exercise_order" example:"1"`
	ExerciseName  string `json:"exercise_name,omitempty" example:"Bench Press"`
	DefaultSets   *int   `json:"default_sets,omitempty" example:"4"`
	DefaultRest   *int   `json:"default_rest,omitempty" example:"90"`
	DefaultReps   *int   `json:"default_reps,omitempty" example:"10"`
}

// AddExerciseToWorkout godoc
// @Summary      Add exercise to workout
// @Description  Assigns an exercise to a workout plan
// @Tags         param_workouts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        workoutId  path      int  true  "Workout ID"
// @Param        request    body      addExerciseToWorkoutRequest  true  "payload"
// @Success      201        {object}  workoutExerciseResponse
// @Failure      400        {object}  presenter.APIError
// @Failure      404        {object}  presenter.APIError
// @Failure      409        {object}  presenter.APIError
// @Failure      500        {object}  presenter.APIError
// @Router       /v1/param_workouts/{workoutId}/exercises [post]
func (h *ParamWorkoutExerciseHandler) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := parseIDParam(chi.URLParam(r, "workoutId"))
	if err != nil {
		status, apiErr := presenter.StatusAndError(param_workout_exercise.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}

	var body struct {
		ExerciseID    int64 `json:"exercise_id"`
		ExerciseOrder int   `json:"exercise_order"`
		DefaultSets   *int  `json:"default_sets"`
		DefaultRest   *int  `json:"default_rest"`
		DefaultReps   *int  `json:"default_reps"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	res, err := h.svc.AddExerciseToWorkout(r.Context(), param_workout_exercise_uc.AddWorkoutExerciseDTO{
		WorkoutID:          workoutID,
		ExerciseID:         body.ExerciseID,
		ExerciseOrder:      body.ExerciseOrder,
		DefaultSets:        body.DefaultSets,
		DefaultRestSeconds: body.DefaultRest,
		DefaultReps:        body.DefaultReps,
	})
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusCreated, workoutExerciseResponse{
		WorkoutID:     res.WorkoutID,
		ExerciseID:    res.ExerciseID,
		ExerciseOrder: res.ExerciseOrder,
		ExerciseName:  res.ExerciseName,
		DefaultSets:   res.DefaultSets,
		DefaultRest:   res.DefaultRestSeconds,
		DefaultReps:   res.DefaultReps,
	})
}

// ListWorkoutExercises godoc
// @Summary      List exercises in workout
// @Description  Retrieves exercises assigned to a workout plan
// @Tags         param_workouts
// @Produce      json
// @Security     BearerAuth
// @Param        workoutId  path      int  true  "Workout ID"
// @Success      200        {array}   workoutExerciseResponse
// @Failure      400        {object}  presenter.APIError
// @Failure      404        {object}  presenter.APIError
// @Failure      500        {object}  presenter.APIError
// @Router       /v1/param_workouts/{workoutId}/exercises [get]
func (h *ParamWorkoutExerciseHandler) ListWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	workoutID, err := parseIDParam(chi.URLParam(r, "workoutId"))
	if err != nil {
		status, apiErr := presenter.StatusAndError(param_workout_exercise.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}

	assignments, err := h.svc.ListWorkoutExercises(r.Context(), workoutID)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	response := make([]workoutExerciseResponse, 0, len(assignments))
	for _, assignment := range assignments {
		response = append(response, workoutExerciseResponse{
			WorkoutID:     assignment.WorkoutID,
			ExerciseID:    assignment.ExerciseID,
			ExerciseOrder: assignment.ExerciseOrder,
			ExerciseName:  assignment.ExerciseName,
			DefaultSets:   assignment.DefaultSets,
			DefaultRest:   assignment.DefaultRestSeconds,
			DefaultReps:   assignment.DefaultReps,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// RemoveExerciseFromWorkout godoc
// @Summary      Remove exercise from workout
// @Description  Removes an exercise assignment from a workout plan
// @Tags         param_workouts
// @Security     BearerAuth
// @Param        workoutId   path  int  true  "Workout ID"
// @Param        exerciseId  path  int  true  "Exercise ID"
// @Success      204
// @Failure      400  {object}  presenter.APIError
// @Failure      404  {object}  presenter.APIError
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_workouts/{workoutId}/exercises/{exerciseId} [delete]
func (h *ParamWorkoutExerciseHandler) RemoveExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := parseIDParam(chi.URLParam(r, "workoutId"))
	if err != nil {
		status, apiErr := presenter.StatusAndError(param_workout_exercise.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}
	exerciseID, err := parseIDParam(chi.URLParam(r, "exerciseId"))
	if err != nil {
		status, apiErr := presenter.StatusAndError(param_workout_exercise.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}

	if err := h.svc.RemoveExerciseFromWorkout(r.Context(), workoutID, exerciseID); err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, param_workout_exercise.ErrInvalidInput
	}
	return id, nil
}
