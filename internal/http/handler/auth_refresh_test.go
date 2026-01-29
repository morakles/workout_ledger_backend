package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	authDomain "workout_ledger/domain/auth"
	token "workout_ledger/internal/auth"
	authUC "workout_ledger/internal/usecase/auth"
)

type fakeAuthRepo struct {
	refreshToken authDomain.RefreshToken
	user         authDomain.User
	rotateCalled bool
}

func (f *fakeAuthRepo) UpsertGoogleUser(ctx context.Context, profile authUC.GoogleProfile) (authDomain.User, error) {
	return authDomain.User{}, nil
}

func (f *fakeAuthRepo) CreateRefreshToken(ctx context.Context, token authDomain.RefreshToken) error {
	return nil
}

func (f *fakeAuthRepo) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (authDomain.RefreshToken, bool, error) {
	if tokenHash != f.refreshToken.TokenHash {
		return authDomain.RefreshToken{}, false, nil
	}
	return f.refreshToken, true, nil
}

func (f *fakeAuthRepo) RotateRefreshToken(ctx context.Context, tokenID int64, revokedAt time.Time, newToken authDomain.RefreshToken) error {
	f.rotateCalled = true
	return nil
}

func (f *fakeAuthRepo) GetUserByID(ctx context.Context, userID int64) (authDomain.User, error) {
	return f.user, nil
}

func TestRefreshTokenHappyPath(t *testing.T) {
	manager, err := token.NewManager("secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	refreshToken := "refresh-token"
	repo := &fakeAuthRepo{
		refreshToken: authDomain.RefreshToken{
			ID:        1,
			UserID:    42,
			TokenHash: manager.HashRefreshToken(refreshToken),
			IssuedAt:  time.Now().UTC().Add(-time.Minute),
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
			CreatedAt: time.Now().UTC().Add(-time.Minute),
		},
		user: authDomain.User{
			ID:    42,
			Email: "user@example.com",
		},
	}
	svc := authUC.NewAuthService(repo, manager)
	handler := NewAuthHandler(svc, nil)

	payload, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.RefreshToken(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected json body, got error: %v", err)
	}
	if body["access_token"] == "" {
		t.Fatalf("expected access token")
	}
	if body["refresh_token"] == "" {
		t.Fatalf("expected refresh token")
	}
	if body["token_type"] != "Bearer" {
		t.Fatalf("expected Bearer token type")
	}
	if !repo.rotateCalled {
		t.Fatalf("expected refresh token rotation")
	}
}

func TestRefreshTokenExpired(t *testing.T) {
	manager, err := token.NewManager("secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	refreshToken := "expired-token"
	repo := &fakeAuthRepo{
		refreshToken: authDomain.RefreshToken{
			ID:        1,
			UserID:    42,
			TokenHash: manager.HashRefreshToken(refreshToken),
			IssuedAt:  time.Now().UTC().Add(-48 * time.Hour),
			ExpiresAt: time.Now().UTC().Add(-24 * time.Hour),
			CreatedAt: time.Now().UTC().Add(-48 * time.Hour),
		},
		user: authDomain.User{
			ID:    42,
			Email: "user@example.com",
		},
	}
	svc := authUC.NewAuthService(repo, manager)
	handler := NewAuthHandler(svc, nil)

	payload, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	handler.RefreshToken(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected json body, got error: %v", err)
	}
	if body["code"] != "AUTH_REFRESH_EXPIRED" {
		t.Fatalf("expected AUTH_REFRESH_EXPIRED, got %v", body["code"])
	}
}
