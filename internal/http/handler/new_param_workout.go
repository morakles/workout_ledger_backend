package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	authDomain "workout_ledger/domain/auth"
	param_workout "workout_ledger/domain/param_workout"
	authmiddleware "workout_ledger/internal/http/middleware"
	"workout_ledger/internal/http/presenter"
	param_workout_uc "workout_ledger/internal/usecase/param_workout"

	"github.com/go-chi/chi/v5"
)

type ParamWorkoutHandler struct {
	svc param_workout_uc.ParamWorkoutService
}

func NewParamWorkoutHandler(svc *param_workout_uc.ParamWorkoutService) *ParamWorkoutHandler {
	return &ParamWorkoutHandler{svc: *svc}
}

// swagger:parameters createParamWorkout
type createParamWorkoutRequest struct {
	Name                   string `json:"name" example:"Push Day"`
	NumberOfSets           int    `json:"number_of_sets" example:"4"`
	RestBetweenSetsSeconds int    `json:"rest_between_sets_seconds" example:"90"`
}

// swagger:response createParamWorkoutResponse
type createParamWorkoutResponse struct {
	Body paramWorkoutResponse `json:"body"`
}

// swagger:response getParamWorkoutsResponse
type getParamWorkoutsResponse struct {
	Body []paramWorkoutResponse `json:"body"`
}

// swagger:parameters getParamWorkoutByID
type getParamWorkoutByIDRequest struct {
	// ID param workout ID
	//
	// in: path
	ID int64 `json:"id" example:"1"`
}

// swagger:response getParamWorkoutByIDResponse
type getParamWorkoutByIDResponse struct {
	Body paramWorkoutResponse `json:"body"`
}

type paramWorkoutResponse struct {
	ID                     int64  `json:"id" example:"1"`
	Name                   string `json:"name" example:"Push Day"`
	NumberOfSets           int    `json:"number_of_sets" example:"4"`
	RestBetweenSetsSeconds int    `json:"rest_between_sets_seconds" example:"90"`
}

// CreateParamWorkout godoc
// @Summary      Create param workout
// @Description  Creates a new param workout
// @Tags         param_workouts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      createParamWorkoutRequest  true  "payload"
// @Success      201      {object}  paramWorkoutResponse
// @Failure      400      {object}  presenter.APIError
// @Failure      409      {object}  presenter.APIError
// @Failure      500      {object}  presenter.APIError
// @Router       /v1/param_workouts [post]
func (h *ParamWorkoutHandler) CreateParamWorkout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name                   string `json:"name"`
		NumberOfSets           int    `json:"number_of_sets"`
		RestBetweenSetsSeconds int    `json:"rest_between_sets_seconds"`
	}

	_ = json.NewDecoder(r.Body).Decode(&body)
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		status, apiErr := presenter.StatusAndError(authDomain.ErrUnauthorized)
		writeJSON(w, status, apiErr)
		return
	}

	res, err := h.svc.CreateParamWorkout(r.Context(), param_workout_uc.CreateParamWorkoutDTO{
		UserID:                 userID,
		Name:                   body.Name,
		NumberOfSets:           body.NumberOfSets,
		RestBetweenSetsSeconds: body.RestBetweenSetsSeconds,
	})
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusCreated, paramWorkoutResponse{
		ID:                     res.ID,
		Name:                   res.Name,
		NumberOfSets:           res.NumberOfSets,
		RestBetweenSetsSeconds: res.RestBetweenSetsSeconds,
	})
}

// GetParamWorkouts godoc
// @Summary      Get param workouts
// @Description  Retrieves all param workouts
// @Tags         param_workouts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   paramWorkoutResponse
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_workouts [get]
func (h *ParamWorkoutHandler) GetParamWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		status, apiErr := presenter.StatusAndError(authDomain.ErrUnauthorized)
		writeJSON(w, status, apiErr)
		return
	}
	workouts, err := h.svc.GetParamWorkouts(r.Context(), userID)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	response := make([]paramWorkoutResponse, 0, len(workouts))
	for _, workout := range workouts {
		response = append(response, paramWorkoutResponse{
			ID:                     workout.ID,
			Name:                   workout.Name,
			NumberOfSets:           workout.NumberOfSets,
			RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// GetParamWorkoutByID godoc
// @Summary      Get param workout by ID
// @Description  Retrieves a param workout by ID
// @Tags         param_workouts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Param workout ID"
// @Success      200  {object}  paramWorkoutResponse
// @Failure      400  {object}  presenter.APIError
// @Failure      404  {object}  presenter.APIError
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_workouts/{id} [get]
func (h *ParamWorkoutHandler) GetParamWorkoutByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		status, apiErr := presenter.StatusAndError(param_workout.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}

	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		status, apiErr := presenter.StatusAndError(authDomain.ErrUnauthorized)
		writeJSON(w, status, apiErr)
		return
	}

	workout, err := h.svc.GetParamWorkoutByID(r.Context(), userID, id)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusOK, paramWorkoutResponse{
		ID:                     workout.ID,
		Name:                   workout.Name,
		NumberOfSets:           workout.NumberOfSets,
		RestBetweenSetsSeconds: workout.RestBetweenSetsSeconds,
	})
}
