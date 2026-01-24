package auth

import "time"

type GoogleProfile struct {
	ProviderUserID string
	Email          string
	EmailVerified  bool
	DisplayName    string
	AvatarURL      string
}

type UserDTO struct {
	ID              int64
	Email           string
	DisplayName     string
	AvatarURL       string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
}
