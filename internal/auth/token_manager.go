package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	now        func() time.Time
}

type AccessTokenClaims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Provider  string `json:"provider"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type RefreshToken struct {
	Token     string
	Hash      string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

func NewManager(secret string, accessTTL, refreshTTL time.Duration) (*Manager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if accessTTL <= 0 || refreshTTL <= 0 {
		return nil, errors.New("token ttl must be positive")
	}
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		now:        time.Now,
	}, nil
}

func (m *Manager) CreateAccessToken(userID int64, email, provider string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(m.accessTTL)
	claims := AccessTokenClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Email:     email,
		Provider:  provider,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}
	signed, err := m.signJWT(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func (m *Manager) CreateRefreshToken(now time.Time) (RefreshToken, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return RefreshToken{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(token))
	return RefreshToken{
		Token:     token,
		Hash:      hex.EncodeToString(hash[:]),
		IssuedAt:  now,
		ExpiresAt: now.Add(m.refreshTTL),
	}, nil
}

func (m *Manager) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (m *Manager) ParseAccessToken(token string) (*AccessTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrTokenInvalid
	}
	signingInput := strings.Join(parts[:2], ".")
	expected := m.sign(signingInput)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, ErrTokenInvalid
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrTokenInvalid
	}
	var claims AccessTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrTokenInvalid
	}
	if m.now().UTC().Unix() >= claims.ExpiresAt {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}

func (m *Manager) signJWT(claims AccessTokenClaims) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + payloadEncoded
	signature := m.sign(signingInput)
	return signingInput + "." + signature, nil
}

func (m *Manager) sign(input string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
