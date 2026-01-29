package handler

import (
	"encoding/json"
	"net/http"
	"workout_ledger/internal/http/presenter"
)

// swagger:parameters refreshTokenRequest
type refreshTokenRequest struct {
	// Refresh token issued during login
	//
	// in: body
	RefreshToken string `json:"refresh_token" example:"refresh_token_value"`
}

// swagger:response refreshTokenResponse
type refreshTokenResponse struct {
	User         authUserResponse `json:"user"`
	AccessToken  string           `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string           `json:"refresh_token" example:"refresh_token_value"`
	TokenType    string           `json:"token_type" example:"Bearer"`
	ExpiresIn    int64            `json:"expires_in" example:"900"`
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchanges a refresh token for a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshTokenRequest  true  "payload"
// @Success      200      {object}  refreshTokenResponse
// @Failure      400      {object}  presenter.APIError
// @Failure      401      {object}  presenter.APIError
// @Failure      500      {object}  presenter.APIError
// @Router       /v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, presenter.APIError{Code: "INVALID_INPUT", Message: "Invalid input"})
		return
	}
	tokens, user, err := h.svc.RefreshTokens(r.Context(), body.RefreshToken)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}
	writeJSON(w, http.StatusOK, refreshTokenResponse{
		User: authUserResponse{
			ID:              user.ID,
			Email:           user.Email,
			DisplayName:     user.DisplayName,
			AvatarURL:       user.AvatarURL,
			EmailVerifiedAt: user.EmailVerifiedAt,
			LastLoginAt:     user.LastLoginAt,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	})
}
