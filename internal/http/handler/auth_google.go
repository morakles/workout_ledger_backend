package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
	authDomain "workout_ledger/domain/auth"
	"workout_ledger/internal/http/presenter"
	authUC "workout_ledger/internal/usecase/auth"
)

const oauthStateCookie = "oauth_state"

type AuthHandler struct {
	svc    *authUC.AuthService
	config *GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AllowedDomain string
}

func NewAuthHandler(svc *authUC.AuthService, cfg *GoogleOAuthConfig) *AuthHandler {
	return &AuthHandler{
		svc:    svc,
		config: cfg,
	}
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	IDToken     string `json:"id_token"`
}

type authUserResponse struct {
	ID              int64      `json:"id" example:"1"`
	Email           string     `json:"email" example:"user@gmail.com"`
	DisplayName     string     `json:"display_name,omitempty" example:"User Name"`
	AvatarURL       string     `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" format:"date-time"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" format:"date-time"`
}

// swagger:response googleLoginResponse
type googleLoginResponse struct {
	User         authUserResponse `json:"user"`
	AccessToken  string           `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string           `json:"refresh_token" example:"refresh_token_value"`
	TokenType    string           `json:"token_type" example:"Bearer"`
	ExpiresIn    int64            `json:"expires_in" example:"900"`
}

// StartGoogleLogin godoc
// @Summary      Start Google login
// @Description  Redirects the user to Google OAuth login
// @Tags         auth
// @Success      307  {string}  string  "Temporary Redirect"
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/auth/google/login [get]
func (h *AuthHandler) StartGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.config == nil {
		writeJSON(w, http.StatusInternalServerError, presenter.APIError{Code: "OAUTH_DISABLED", Message: "Google OAuth is not configured"})
		return
	}

	state, err := generateState()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, presenter.APIError{Code: "STATE_GENERATION_FAILED", Message: "Failed to generate state"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
		Expires:  time.Now().Add(5 * time.Minute),
	})

	authURL := h.authCodeURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback godoc
// @Summary      Handle Google login callback
// @Description  Exchanges the OAuth code for user info and logs in the user
// @Tags         auth
// @Produce      json
// @Param        code   query     string true  "OAuth code"
// @Param        state  query     string true  "OAuth state"
// @Success      200  {object}  googleLoginResponse
// @Failure      400  {object}  presenter.APIError
// @Failure      401  {object}  presenter.APIError
// @Failure      500  {object}  presenter.APIError
// @Router       /v1/auth/google/callback [get]
func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.config == nil {
		writeJSON(w, http.StatusInternalServerError, presenter.APIError{Code: "OAUTH_DISABLED", Message: "Google OAuth is not configured"})
		return
	}

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if code == "" || state == "" {
		writeJSON(w, http.StatusBadRequest, presenter.APIError{Code: "INVALID_INPUT", Message: "Missing code or state"})
		return
	}

	if err := h.validateState(r, state); err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	accessToken, err := h.exchangeCodeForToken(r.Context(), code)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, presenter.APIError{Code: "OAUTH_EXCHANGE_FAILED", Message: "Failed to exchange OAuth code"})
		return
	}

	userInfo, err := h.fetchGoogleUserInfo(r.Context(), accessToken)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, presenter.APIError{Code: "OAUTH_USERINFO_FAILED", Message: "Failed to fetch user info"})
		return
	}

	if userInfo.Email == "" || userInfo.ID == "" {
		writeJSON(w, http.StatusBadRequest, presenter.APIError{Code: "INVALID_PROFILE", Message: "Google profile is missing required fields"})
		return
	}

	if h.config.AllowedDomain != "" && !strings.HasSuffix(strings.ToLower(userInfo.Email), "@"+strings.ToLower(h.config.AllowedDomain)) {
		writeJSON(w, http.StatusUnauthorized, presenter.APIError{Code: "DOMAIN_NOT_ALLOWED", Message: "Email domain is not allowed"})
		return
	}

	user, err := h.svc.LoginWithGoogle(r.Context(), authUC.GoogleProfile{
		ProviderUserID: userInfo.ID,
		Email:          userInfo.Email,
		EmailVerified:  userInfo.VerifiedEmail,
		DisplayName:    userInfo.Name,
		AvatarURL:      userInfo.Picture,
	})
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	tokens, err := h.svc.IssueTokens(r.Context(), user, "google")
	if err != nil {
		status, apiErr := presenter.StatusAndError(err)
		writeJSON(w, status, apiErr)
		return
	}

	writeJSON(w, http.StatusOK, googleLoginResponse{
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

func (h *AuthHandler) authCodeURL(state string) string {
	values := url.Values{}
	values.Set("client_id", h.config.ClientID)
	values.Set("redirect_uri", h.config.RedirectURL)
	values.Set("response_type", "code")
	values.Set("scope", strings.Join([]string{"openid", "email", "profile"}, " "))
	values.Set("state", state)
	values.Set("access_type", "offline")
	return "https://accounts.google.com/o/oauth2/v2/auth?" + values.Encode()
}

func (h *AuthHandler) exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", h.config.ClientID)
	form.Set("client_secret", h.config.ClientSecret)
	form.Set("redirect_uri", h.config.RedirectURL)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", authDomain.ErrUnauthorized
	}

	var tokenResponse googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}
	if tokenResponse.AccessToken == "" {
		return "", authDomain.ErrUnauthorized
	}
	return tokenResponse.AccessToken, nil
}

func (h *AuthHandler) fetchGoogleUserInfo(ctx context.Context, accessToken string) (googleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return googleUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return googleUserInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return googleUserInfo{}, authDomain.ErrUnauthorized
	}
	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return googleUserInfo{}, err
	}
	return info, nil
}

func (h *AuthHandler) validateState(r *http.Request, state string) error {
	cookie, err := r.Cookie(oauthStateCookie)
	if err != nil || cookie.Value == "" {
		return authDomain.ErrStateMismatch
	}
	if cookie.Value != state {
		return authDomain.ErrStateMismatch
	}
	return nil
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
