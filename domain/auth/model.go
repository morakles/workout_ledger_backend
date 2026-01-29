package auth

import "time"

type User struct {
	ID              int64
	Email           string
	DisplayName     string
	AvatarURL       string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	AuthProvider    string
	PasswordHash    *string
}

const (
	AuthProviderLocal       = "LOCAL"
	AuthProviderGoogle      = "GOOGLE"
	AuthProviderLocalGoogle = "LOCAL_GOOGLE"
)
