package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	authDomain "workout_ledger/domain/auth"
	token "workout_ledger/internal/auth"
	"workout_ledger/internal/http/presenter"
)

type contextKey string

const userIDKey contextKey = "user_id"

func UserIDFromContext(ctx context.Context) (int64, bool) {
	value := ctx.Value(userIDKey)
	id, ok := value.(int64)
	return id, ok
}

func JWTAuth(tokenManager *token.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				respondAuthError(w, authDomain.ErrMissingToken)
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				respondAuthError(w, authDomain.ErrInvalidToken)
				return
			}
			rawToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if rawToken == "" {
				respondAuthError(w, authDomain.ErrMissingToken)
				return
			}
			claims, err := tokenManager.ParseAccessToken(rawToken)
			if err != nil {
				if errors.Is(err, token.ErrTokenExpired) {
					respondAuthError(w, authDomain.ErrTokenExpired)
					return
				}
				respondAuthError(w, authDomain.ErrInvalidToken)
				return
			}
			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil || userID <= 0 {
				respondAuthError(w, authDomain.ErrInvalidToken)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondAuthError(w http.ResponseWriter, err error) {
	status, apiErr := presenter.StatusAndError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiErr)
}
