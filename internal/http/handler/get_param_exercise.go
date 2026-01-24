package handler

// swagger:response listParamExercisesResponse
type listParamExercisesResponse struct {
	Items []struct {
		ID      int64  `json:"id" example:"1"`
		Name    string `json:"name" example:"Push-up"`
		IconUrl string `json:"icon_url" example:"https://example.com/icon.png"`
	} `json:"items"`
}
