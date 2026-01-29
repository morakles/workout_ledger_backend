package auth

import "time"

type User struct {
	ID              int64
	Email           string
	DisplayName     string
	AvatarURL       string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
}
