package auth

import (
	"context"
	"net/mail"
	"strings"
	"time"
	authDomain "workout_ledger/domain/auth"
	token "workout_ledger/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	UpsertGoogleUser(ctx context.Context, profile GoogleProfile) (authDomain.User, error)
	CreateLocalUser(ctx context.Context, email, passwordHash string, now time.Time) (authDomain.User, error)
	FindUserByEmail(ctx context.Context, email string) (authDomain.User, bool, error)
	UpdateLastLogin(ctx context.Context, userID int64, now time.Time) error
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
		AuthProvider:    user.AuthProvider,
	}, nil
}

func (s *AuthService) RegisterLocalUser(ctx context.Context, email, password string) (UserDTO, error) {
	cleanEmail := strings.TrimSpace(email)
	if !isValidEmail(cleanEmail) {
		return UserDTO{}, authDomain.ErrInvalidInput
	}
	if len(password) < minPasswordLength {
		return UserDTO{}, authDomain.ErrInvalidInput
	}
	_, found, err := s.repo.FindUserByEmail(ctx, cleanEmail)
	if err != nil {
		return UserDTO{}, err
	}
	if found {
		return UserDTO{}, authDomain.ErrConflict
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserDTO{}, err
	}
	now := time.Now().UTC()
	user, err := s.repo.CreateLocalUser(ctx, cleanEmail, string(hash), now)
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
		AuthProvider:    user.AuthProvider,
	}, nil
}

func (s *AuthService) LoginWithEmail(ctx context.Context, email, password string) (UserDTO, error) {
	cleanEmail := strings.TrimSpace(email)
	if !isValidEmail(cleanEmail) || strings.TrimSpace(password) == "" {
		return UserDTO{}, authDomain.ErrInvalidInput
	}
	user, found, err := s.repo.FindUserByEmail(ctx, cleanEmail)
	if err != nil {
		return UserDTO{}, err
	}
	if !found || user.PasswordHash == nil {
		return UserDTO{}, authDomain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return UserDTO{}, authDomain.ErrUnauthorized
	}
	now := time.Now().UTC()
	if err := s.repo.UpdateLastLogin(ctx, user.ID, now); err != nil {
		return UserDTO{}, err
	}
	return UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		DisplayName:     user.DisplayName,
		AvatarURL:       user.AvatarURL,
		EmailVerifiedAt: user.EmailVerifiedAt,
		LastLoginAt:     &now,
		AuthProvider:    user.AuthProvider,
	}, nil
}

func (s *AuthService) IssueTokens(ctx context.Context, user UserDTO) (TokenPair, error) {
	if s.tokenManager == nil {
		return TokenPair{}, authDomain.ErrUnauthorized
	}
	now := time.Now().UTC()
	accessToken, accessExpiresAt, err := s.tokenManager.CreateAccessToken(user.ID, user.Email, tokenProvider(user.AuthProvider), now)
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
		AuthProvider:    user.AuthProvider,
	}
	accessToken, accessExpiresAt, err := s.tokenManager.CreateAccessToken(user.ID, user.Email, tokenProvider(user.AuthProvider), now)
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

const minPasswordLength = 8

func isValidEmail(value string) bool {
	if value == "" {
		return false
	}
	parsed, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}
	return parsed.Address == value
}

func tokenProvider(provider string) string {
	switch strings.ToUpper(strings.TrimSpace(provider)) {
	case authDomain.AuthProviderLocal:
		return "local"
	case authDomain.AuthProviderLocalGoogle:
		return "local_google"
	case authDomain.AuthProviderGoogle:
		return "google"
	default:
		return "google"
	}
}
