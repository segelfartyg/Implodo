package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Google OAuth2 credentials — create these at https://console.cloud.google.com/
	// Application type: Web application
	// Authorized redirect URI must include GoogleRedirectURL below
	GoogleClientID     string
	GoogleClientSecret string

	// Where Google sends the user after login — must match what you registered in Google Cloud Console
	GoogleRedirectURL string

	// Hardcoded list of allowed Google user IDs ("sub" field from Google's userinfo endpoint)
	// Find your Google user ID by logging in and checking the /api/me response
	AllowedGoogleIDs []string

	// Secret used to sign JWTs — replace with a long random string in production
	// Generate one: openssl rand -hex 32
	JWTSecret string

	// How long issued JWTs remain valid
	JWTDuration time.Duration

	// Deep link URI scheme the mobile app registers — used to signal auth completion to the app
	// e.g. "myapp://auth/complete" — the app intercepts this and knows to call /auth/token
	AppDeepLinkURI string

	Port       string
	BucketName string
}

func Load() *Config {
	hours, _ := strconv.Atoi(getEnv("JWT_DURATION_HOURS", "24"))

	rawIDs := getEnv("ALLOWED_GOOGLE_IDS", "")
	var allowedIDs []string
	for _, id := range strings.Split(rawIDs, ",") {
		if id = strings.TrimSpace(id); id != "" {
			allowedIDs = append(allowedIDs, id)
		}
	}

	return &Config{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		AllowedGoogleIDs:   allowedIDs,
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTDuration:        time.Duration(hours) * time.Hour,
		AppDeepLinkURI:     getEnv("APP_DEEP_LINK_URI", "implodo://auth/complete"),
		Port:               getEnv("PORT", "8080"),
		BucketName:         getEnv("BUCKET_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
