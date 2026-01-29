package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddress string
	DSN         string
	GoogleOAuth *GoogleOAuthConfig
	JWT         JWTConfig
}

type GoogleOAuthConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AllowedDomain string
}

type JWTConfig struct {
	Secret      string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
	AccessMins  int
	RefreshDays int
}

func Load() (Config, error) {
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	postgresdsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if postgresdsn == "" {
		return Config{}, fmt.Errorf("environment variable POSTGRES_DSN is not set")
	}
	googleOAuth, err := loadGoogleOAuthConfig()
	if err != nil {
		return Config{}, err
	}
	jwtConfig, err := loadJWTConfig()
	if err != nil {
		return Config{}, err
	}
	return Config{
		HTTPAddress: httpAddr,
		DSN:         postgresdsn,
		GoogleOAuth: googleOAuth,
		JWT:         jwtConfig,
	}, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func loadGoogleOAuthConfig() (*GoogleOAuthConfig, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	redirectURL := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))
	allowedDomain := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_ALLOWED_DOMAIN"))

	if clientID == "" && clientSecret == "" && redirectURL == "" {
		return nil, nil
	}
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, fmt.Errorf("google oauth requires GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET, and GOOGLE_OAUTH_REDIRECT_URL")
	}

	return &GoogleOAuthConfig{
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		RedirectURL:   redirectURL,
		AllowedDomain: allowedDomain,
	}, nil
}

func loadJWTConfig() (JWTConfig, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return JWTConfig{}, fmt.Errorf("environment variable JWT_SECRET is not set")
	}
	accessMinutes, err := getEnvInt("JWT_ACCESS_TTL_MINUTES", 15)
	if err != nil {
		return JWTConfig{}, err
	}
	refreshDays, err := getEnvInt("JWT_REFRESH_TTL_DAYS", 30)
	if err != nil {
		return JWTConfig{}, err
	}
	return JWTConfig{
		Secret:      secret,
		AccessTTL:   time.Duration(accessMinutes) * time.Minute,
		RefreshTTL:  time.Duration(refreshDays) * 24 * time.Hour,
		AccessMins:  accessMinutes,
		RefreshDays: refreshDays,
	}, nil
}

func getEnvInt(key string, def int) (int, error) {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("environment variable %s must be an integer", key)
		}
		if parsed <= 0 {
			return 0, fmt.Errorf("environment variable %s must be positive", key)
		}
		return parsed, nil
	}
	return def, nil
}
