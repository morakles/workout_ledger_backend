package auth

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	IssuedAt  time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
