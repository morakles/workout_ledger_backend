package auth

import (
	"context"
	"strings"
	authDomain "workout_ledger/domain/auth"
)

type AuthRepository interface {
	UpsertGoogleUser(ctx context.Context, profile GoogleProfile) (authDomain.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, profile GoogleProfile) (UserDTO, error) {
	if strings.TrimSpace(profile.ProviderUserID) == "" || strings.TrimSpace(profile.Email) == "" {
		return UserDTO{}, authDomain.ErrInvalidInput
	}

	user, err := s.repo.UpsertGoogleUser(ctx, profile)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		DisplayName:     user.DisplayName,
		AvatarURL:       user.AvatarURL,
		EmailVerifiedAt: user.EmailVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
	}, nil
}
