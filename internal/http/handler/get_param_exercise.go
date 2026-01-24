package handler

import (
	"net/http"
	"workout_ledger/internal/http/presenter"
	param_exercise_uc "workout_ledger/internal/usecase/param_exercise"
)

// swagger:response listParamExercisesResponse
type listParamExercisesResponse struct {
	Items []struct {
		ID      int64  `json:"id" example:"1"`
		Name    string `json:"name" example:"Push-up"`
		IconUrl string `json:"icon_url" example:"https://example.com/icon.png"`
	} `json:"items"`
}

// ListParamExercises godoc
// @Summary      List param exercises
// @Description  Returns all param exercises
// @Tags         param_exercises
// @Produce      json
// @Success      200  {object}  listParamExercisesResponse
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/param_exercises [get]
func (h *ParamExerciseHandler) ListParamExercises(w http.ResponseWriter, r *http.Request) {
	// Replace this call with whatever your usecase/service exposes.
	items, err := h.svc.ListParamExercises(r.Context(), param_exercise_uc.ListParamExercisesDTO{})
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, map[string]any{
			"id":       it.ID,
			"name":     it.Name,
			"icon_url": it.IconUrl,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": out,
	})
}
