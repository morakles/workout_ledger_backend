package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workout_ledger/domain/param_exercise"
	"workout_ledger/internal/http/presenter"
	param_exercise_uc "workout_ledger/internal/usecase/param_exercise"

	"github.com/go-chi/chi/v5"
)

type ParamExerciseHandler struct {
	svc param_exercise_uc.ParamExerciseService
}

func NewParamExerciseHandler(svc *param_exercise_uc.ParamExerciseService) *ParamExerciseHandler {
	return &ParamExerciseHandler{svc: *svc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// swagger:parameters createParamExercise
type createParamExerciseRequest struct {
	Name    string `json:"name" example:"Push-up"`
	IconUrl string `json:"icon_url" example:"https://example.com/icon.png"`
}

// swagger:response createParamExerciseResponse
type createParamExerciseResponse struct {
	ID int64 `json:"id" example:"1"`
}

// swagger:response getParamExercisesResponse
type getParamExercisesResponse struct {
	// in: body
	Body []paramExerciseResponse `json:"body"`
}

// swagger:parameters getParamExerciseByID
type getParamExerciseByIDRequest struct {
	// ID param exercise ID
	//
	// in: path
	ID int64 `json:"id" example:"1"`
}

// swagger:response getParamExerciseByIDResponse
type getParamExerciseByIDResponse struct {
	// in: body
	Body paramExerciseResponse `json:"body"`
}

type paramExerciseResponse struct {
	ID      int64  `json:"id" example:"1"`
	Name    string `json:"name" example:"Push-up"`
	IconUrl string `json:"icon_url" example:"https://example.com/icon.png"`
}

// CreateParamExercise godoc
// @Summary      Create param exercise
// @Description  Creates a new param exercise
// @Tags         param_exercises
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      createParamExerciseRequest  true  "payload"
// @Success      201      {object}  createParamExerciseResponse
// @Failure      400      {object}  presenter.APIError
// @Failure      409      {object}  presenter.APIError
// @Failure      500      {object}  presenter.APIError
// @Router       /v1/param_exercises [post]
func (h *ParamExerciseHandler) CreateParamExercise(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		IconUrl string `json:"icon_url"`
	}

	_ = json.NewDecoder(r.Body).Decode(&body)

	res, err := h.svc.CreateParamExercise(r.Context(), param_exercise_uc.CreateParamExerciseDTO{
		Name:    body.Name,
		IconUrl: body.IconUrl,
	})
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id": res.ID,
	})
}

// GetParamExercises godoc
// @Summary      Get param exercises
// @Description  Retrieves all param exercises
// @Tags         param_exercises
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   paramExerciseResponse
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_exercises [get]
func (h *ParamExerciseHandler) GetParamExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.svc.GetParamExercises(r.Context())
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	response := make([]paramExerciseResponse, 0, len(exercises))
	for _, exercise := range exercises {
		response = append(response, paramExerciseResponse{
			ID:      exercise.ID,
			Name:    exercise.Name,
			IconUrl: exercise.IconUrl,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// GetParamExerciseByID godoc
// @Summary      Get param exercise by ID
// @Description  Retrieves a param exercise by ID
// @Tags         param_exercises
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Param exercise ID"
// @Success      200  {object}  paramExerciseResponse
// @Failure      400  {object}  presenter.APIError
// @Failure      404  {object}  presenter.APIError
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_exercises/{id} [get]
func (h *ParamExerciseHandler) GetParamExerciseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		status, apiErr := presenter.StatusAndError(param_exercise.ErrInvalidInput)
		writeJSON(w, status, apiErr)
		return
	}

	exercise, err := h.svc.GetParamExerciseByID(r.Context(), id)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusOK, paramExerciseResponse{
		ID:      exercise.ID,
		Name:    exercise.Name,
		IconUrl: exercise.IconUrl,
	})
}
