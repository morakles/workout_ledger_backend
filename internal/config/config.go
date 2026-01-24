package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPAddress string
	DSN         string
	GoogleOAuth *GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AllowedDomain string
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
	return Config{
		HTTPAddress: httpAddr,
		DSN:         postgresdsn,
		GoogleOAuth: googleOAuth,
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
