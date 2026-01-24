package handler

import (
	"encoding/json"
	"net/http"
	"workout_ledger/internal/http/presenter"
	param_exercise_uc "workout_ledger/internal/usecase/param_exercise"
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

// CreateParamExercise godoc
// @Summary      Create param exercise
// @Description  Creates a new param exercise
// @Tags         param_exercises
// @Accept       json
// @Produce      json
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
