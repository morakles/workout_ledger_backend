package auth

import (
	"context"
	"strings"
	"time"
	authDomain "workout_ledger/domain/auth"
	token "workout_ledger/internal/auth"
)

type AuthRepository interface {
	UpsertGoogleUser(ctx context.Context, profile GoogleProfile) (authDomain.User, error)
	CreateRefreshToken(ctx context.Context, token authDomain.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (authDomain.RefreshToken, bool, error)
	RotateRefreshToken(ctx context.Context, tokenID int64, revokedAt time.Time, newToken authDomain.RefreshToken) error
	GetUserByID(ctx context.Context, userID int64) (authDomain.User, error)
}

type AuthService struct {
	repo         AuthRepository
	tokenManager *token.Manager
}

func NewAuthService(repo AuthRepository, tokenManager *token.Manager) *AuthService {
	return &AuthService{repo: repo, tokenManager: tokenManager}
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

func (s *AuthService) IssueTokens(ctx context.Context, user UserDTO, provider string) (TokenPair, error) {
	if s.tokenManager == nil {
		return TokenPair{}, authDomain.ErrUnauthorized
	}
	now := time.Now().UTC()
	accessToken, accessExpiresAt, err := s.tokenManager.CreateAccessToken(user.ID, user.Email, provider, now)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := s.tokenManager.CreateRefreshToken(now)
	if err != nil {
		return TokenPair{}, err
	}
	record := authDomain.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken.Hash,
		IssuedAt:  refreshToken.IssuedAt,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: now,
	}
	if err := s.repo.CreateRefreshToken(ctx, record); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiresAt.Sub(now).Seconds()),
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (TokenPair, UserDTO, error) {
	if s.tokenManager == nil {
		return TokenPair{}, UserDTO{}, authDomain.ErrUnauthorized
	}
	if strings.TrimSpace(refreshToken) == "" {
		return TokenPair{}, UserDTO{}, authDomain.ErrRefreshInvalid
	}
	now := time.Now().UTC()
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)
	stored, found, err := s.repo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return TokenPair{}, UserDTO{}, err
	}
	if !found || stored.RevokedAt != nil {
		return TokenPair{}, UserDTO{}, authDomain.ErrRefreshInvalid
	}
	if now.After(stored.ExpiresAt) {
		return TokenPair{}, UserDTO{}, authDomain.ErrRefreshExpired
	}
	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil {
		return TokenPair{}, UserDTO{}, err
	}
	userDTO := UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		DisplayName:     user.DisplayName,
		AvatarURL:       user.AvatarURL,
		EmailVerifiedAt: user.EmailVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
	}
	accessToken, accessExpiresAt, err := s.tokenManager.CreateAccessToken(user.ID, user.Email, "google", now)
	if err != nil {
		return TokenPair{}, UserDTO{}, err
	}
	newRefreshToken, err := s.tokenManager.CreateRefreshToken(now)
	if err != nil {
		return TokenPair{}, UserDTO{}, err
	}
	newRecord := authDomain.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshToken.Hash,
		IssuedAt:  newRefreshToken.IssuedAt,
		ExpiresAt: newRefreshToken.ExpiresAt,
		CreatedAt: now,
	}
	if err := s.repo.RotateRefreshToken(ctx, stored.ID, now, newRecord); err != nil {
		return TokenPair{}, UserDTO{}, err
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiresAt.Sub(now).Seconds()),
	}, userDTO, nil
}
