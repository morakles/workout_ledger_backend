package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	authDomain "workout_ledger/domain/auth"
	authUC "workout_ledger/internal/usecase/auth"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) UpsertGoogleUser(ctx context.Context, profile authUC.GoogleProfile) (authDomain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return authDomain.User{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	providerID, err := ensureAuthProvider(ctx, tx, "google")
	if err != nil {
		return authDomain.User{}, err
	}

	user, found, err := findUserByIdentity(ctx, tx, providerID, profile.ProviderUserID)
	if err != nil {
		return authDomain.User{}, err
	}
	if !found {
		user, found, err = findUserByEmail(ctx, tx, profile.Email)
		if err != nil {
			return authDomain.User{}, err
		}
		if !found {
			user, err = createUser(ctx, tx, profile)
			if err != nil {
				return authDomain.User{}, err
			}
		}
		if err := createAuthIdentity(ctx, tx, providerID, user.ID, profile); err != nil {
			return authDomain.User{}, err
		}
	}

	now := time.Now().UTC()
	user, err = updateUserFromGoogle(ctx, tx, user.ID, profile, now)
	if err != nil {
		return authDomain.User{}, err
	}

	if err = tx.Commit(); err != nil {
		return authDomain.User{}, err
	}
	return user, nil
}

func ensureAuthProvider(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	var id int64
	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO auth_provider (provider_name)
		 VALUES ($1)
		 ON CONFLICT (provider_name) DO UPDATE SET provider_name = EXCLUDED.provider_name
		 RETURNING id`,
		name,
	).Scan(&id)
	return id, err
}

func findUserByIdentity(ctx context.Context, tx *sql.Tx, providerID int64, providerUserID string) (authDomain.User, bool, error) {
	var user authDomain.User
	var displayName sql.NullString
	var avatarURL sql.NullString
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime

	err := tx.QueryRowContext(
		ctx,
		`SELECT u.id, u.email, u.display_name, u.avatar_url, u.email_verified_at, u.last_login_at
		 FROM auth_identities ai
		 JOIN param_users u ON u.id = ai.user_id
		 WHERE ai.auth_provider_id = $1 AND ai.provider_user_id = $2`,
		providerID,
		providerUserID,
	).Scan(&user.ID, &user.Email, &displayName, &avatarURL, &emailVerifiedAt, &lastLoginAt)
	if errors.Is(err, sql.ErrNoRows) {
		return authDomain.User{}, false, nil
	}
	if err != nil {
		return authDomain.User{}, false, err
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if emailVerifiedAt.Valid {
		t := emailVerifiedAt.Time
		user.EmailVerifiedAt = &t
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time
		user.LastLoginAt = &t
	}
	return user, true, nil
}

func findUserByEmail(ctx context.Context, tx *sql.Tx, email string) (authDomain.User, bool, error) {
	var user authDomain.User
	var displayName sql.NullString
	var avatarURL sql.NullString
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime

	err := tx.QueryRowContext(
		ctx,
		`SELECT id, email, display_name, avatar_url, email_verified_at, last_login_at
		 FROM param_users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &displayName, &avatarURL, &emailVerifiedAt, &lastLoginAt)
	if errors.Is(err, sql.ErrNoRows) {
		return authDomain.User{}, false, nil
	}
	if err != nil {
		return authDomain.User{}, false, err
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if emailVerifiedAt.Valid {
		t := emailVerifiedAt.Time
		user.EmailVerifiedAt = &t
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time
		user.LastLoginAt = &t
	}
	return user, true, nil
}

func createUser(ctx context.Context, tx *sql.Tx, profile authUC.GoogleProfile) (authDomain.User, error) {
	var user authDomain.User
	var emailVerifiedAt sql.NullTime
	now := time.Now().UTC()
	if profile.EmailVerified {
		emailVerifiedAt = sql.NullTime{Time: now, Valid: true}
	}
	var displayName sql.NullString
	if profile.DisplayName != "" {
		displayName = sql.NullString{String: profile.DisplayName, Valid: true}
	}
	var avatarURL sql.NullString
	if profile.AvatarURL != "" {
		avatarURL = sql.NullString{String: profile.AvatarURL, Valid: true}
	}

	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO param_users (email, email_verified_at, display_name, avatar_url, last_login_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, email, email_verified_at, display_name, avatar_url, last_login_at`,
		profile.Email,
		emailVerifiedAt,
		displayName,
		avatarURL,
		now,
	).Scan(&user.ID, &user.Email, &emailVerifiedAt, &displayName, &avatarURL, &user.LastLoginAt)
	if err != nil {
		return authDomain.User{}, err
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	return user, nil
}

func createAuthIdentity(ctx context.Context, tx *sql.Tx, providerID, userID int64, profile authUC.GoogleProfile) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO auth_identities (user_id, auth_provider_id, provider_user_id, email_at_provider)
		 VALUES ($1, $2, $3, $4)`,
		userID,
		providerID,
		profile.ProviderUserID,
		profile.Email,
	)
	return err
}

func updateUserFromGoogle(ctx context.Context, tx *sql.Tx, userID int64, profile authUC.GoogleProfile, now time.Time) (authDomain.User, error) {
	var displayName sql.NullString
	if profile.DisplayName != "" {
		displayName = sql.NullString{String: profile.DisplayName, Valid: true}
	}
	var avatarURL sql.NullString
	if profile.AvatarURL != "" {
		avatarURL = sql.NullString{String: profile.AvatarURL, Valid: true}
	}
	var emailVerifiedAt sql.NullTime
	if profile.EmailVerified {
		emailVerifiedAt = sql.NullTime{Time: now, Valid: true}
	}

	var user authDomain.User
	var emailVerifiedAtResult sql.NullTime
	var displayNameResult sql.NullString
	var avatarURLResult sql.NullString
	var lastLoginAt sql.NullTime
	err := tx.QueryRowContext(
		ctx,
		`UPDATE param_users
		 SET email = $1,
		     email_verified_at = COALESCE($2, email_verified_at),
		     display_name = COALESCE($3, display_name),
		     avatar_url = COALESCE($4, avatar_url),
		     updated_at = $5,
		     last_login_at = $5
		 WHERE id = $6
		 RETURNING id, email, email_verified_at, display_name, avatar_url, last_login_at`,
		profile.Email,
		emailVerifiedAt,
		displayName,
		avatarURL,
		now,
		userID,
	).Scan(&user.ID, &user.Email, &emailVerifiedAtResult, &displayNameResult, &avatarURLResult, &lastLoginAt)
	if err != nil {
		return authDomain.User{}, err
	}
	if emailVerifiedAtResult.Valid {
		user.EmailVerifiedAt = &emailVerifiedAtResult.Time
	}
	if displayNameResult.Valid {
		user.DisplayName = displayNameResult.String
	}
	if avatarURLResult.Valid {
		user.AvatarURL = avatarURLResult.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	return user, nil
}

func (r *AuthRepo) CreateRefreshToken(ctx context.Context, token authDomain.RefreshToken) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO auth_refresh_tokens (user_id, token_hash, issued_at, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		token.UserID,
		token.TokenHash,
		token.IssuedAt,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *AuthRepo) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (authDomain.RefreshToken, bool, error) {
	var token authDomain.RefreshToken
	var revokedAt sql.NullTime
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token_hash, issued_at, expires_at, revoked_at, created_at
		 FROM auth_refresh_tokens
		 WHERE token_hash = $1`,
		tokenHash,
	).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.IssuedAt, &token.ExpiresAt, &revokedAt, &token.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return authDomain.RefreshToken{}, false, nil
	}
	if err != nil {
		return authDomain.RefreshToken{}, false, err
	}
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	return token, true, nil
}

func (r *AuthRepo) RotateRefreshToken(ctx context.Context, tokenID int64, revokedAt time.Time, newToken authDomain.RefreshToken) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`UPDATE auth_refresh_tokens
		 SET revoked_at = $1
		 WHERE id = $2 AND revoked_at IS NULL`,
		revokedAt,
		tokenID,
	)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO auth_refresh_tokens (user_id, token_hash, issued_at, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		newToken.UserID,
		newToken.TokenHash,
		newToken.IssuedAt,
		newToken.ExpiresAt,
		newToken.CreatedAt,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *AuthRepo) GetUserByID(ctx context.Context, userID int64) (authDomain.User, error) {
	var user authDomain.User
	var displayName sql.NullString
	var avatarURL sql.NullString
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, display_name, avatar_url, email_verified_at, last_login_at
		 FROM param_users
		 WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Email, &displayName, &avatarURL, &emailVerifiedAt, &lastLoginAt)
	if err != nil {
		return authDomain.User{}, err
	}
	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if emailVerifiedAt.Valid {
		t := emailVerifiedAt.Time
		user.EmailVerifiedAt = &t
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time
		user.LastLoginAt = &t
	}
	return user, nil
}
