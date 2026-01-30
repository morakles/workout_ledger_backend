package handler

import (
	"encoding/json"
	"net/http"
	"workout_ledger/internal/http/presenter"
)

// swagger:parameters registerRequest
type registerRequest struct {
	// in: body
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongpassword"`
}

// swagger:response registerResponse
type registerResponse struct {
	User         authUserResponse `json:"user"`
	AccessToken  string           `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string           `json:"refresh_token" example:"refresh_token_value"`
	TokenType    string           `json:"token_type" example:"Bearer"`
	ExpiresIn    int64            `json:"expires_in" example:"900"`
}

// Register godoc
// @Summary      Register with email and password
// @Description  Creates a local account and returns tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest  true  "payload"
// @Success      201      {object}  registerResponse
// @Failure      400      {object}  presenter.APIError
// @Failure      409      {object}  presenter.APIError
// @Failure      500      {object}  presenter.APIError
// @Router       /v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, presenter.APIError{Code: "INVALID_INPUT", Message: "Invalid input"})
		return
	}
	user, err := h.svc.RegisterLocalUser(r.Context(), body.Email, body.Password)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}
	tokens, err := h.svc.IssueTokens(r.Context(), user)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}
	writeJSON(w, http.StatusCreated, registerResponse{
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

// swagger:parameters loginRequest
type loginRequest struct {
	// in: body
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongpassword"`
}

// swagger:response loginResponse
type loginResponse struct {
	User         authUserResponse `json:"user"`
	AccessToken  string           `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string           `json:"refresh_token" example:"refresh_token_value"`
	TokenType    string           `json:"token_type" example:"Bearer"`
	ExpiresIn    int64            `json:"expires_in" example:"900"`
}

// Login godoc
// @Summary      Login with email and password
// @Description  Authenticates a local account and returns tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "payload"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  presenter.APIError
// @Failure      401      {object}  presenter.APIError
// @Failure      500      {object}  presenter.APIError
// @Router       /v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, presenter.APIError{Code: "INVALID_INPUT", Message: "Invalid input"})
		return
	}
	user, err := h.svc.LoginWithEmail(r.Context(), body.Email, body.Password)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}
	tokens, err := h.svc.IssueTokens(r.Context(), user)
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{
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
