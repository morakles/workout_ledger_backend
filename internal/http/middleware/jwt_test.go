package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	token "workout_ledger/internal/auth"
)

type apiError struct {
	Code string `json:"code"`
}

func TestJWTAuthMissingToken(t *testing.T) {
	manager, err := token.NewManager("secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	handler := JWTAuth(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/v1/param_exercises", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	var body apiError
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected json body, got error: %v", err)
	}
	if body.Code != "AUTH_MISSING_TOKEN" {
		t.Fatalf("expected AUTH_MISSING_TOKEN, got %s", body.Code)
	}
}

func TestJWTAuthInvalidToken(t *testing.T) {
	manager, err := token.NewManager("secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	handler := JWTAuth(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/v1/param_exercises", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	var body apiError
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected json body, got error: %v", err)
	}
	if body.Code != "AUTH_INVALID_TOKEN" {
		t.Fatalf("expected AUTH_INVALID_TOKEN, got %s", body.Code)
	}
}
