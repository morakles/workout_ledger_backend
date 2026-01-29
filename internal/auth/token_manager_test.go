package token

import (
	"testing"
	"time"
)

func TestCreateAndParseAccessToken(t *testing.T) {
	manager, err := NewManager("secret", 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	now := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return now }
	token, exp, err := manager.CreateAccessToken(42, "user@example.com", "google", now)
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}
	if exp.Before(now) {
		t.Fatalf("expected exp after now")
	}
	claims, err := manager.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("expected claims, got error: %v", err)
	}
	if claims.Subject != "42" {
		t.Fatalf("expected subject 42, got %s", claims.Subject)
	}
	if claims.Email != "user@example.com" {
		t.Fatalf("expected email claim, got %s", claims.Email)
	}
	if claims.Provider != "google" {
		t.Fatalf("expected provider claim, got %s", claims.Provider)
	}
}
