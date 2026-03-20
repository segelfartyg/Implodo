package config

import "time"

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

	Port string
}

func Load() *Config {
	return &Config{
		GoogleClientID:     "PLACEHOLDER_GOOGLE_CLIENT_ID.apps.googleusercontent.com",
		GoogleClientSecret: "PLACEHOLDER_GOOGLE_CLIENT_SECRET",
		GoogleRedirectURL:  "http://localhost:8080/auth/google/callback",

		AllowedGoogleIDs: []string{
			"PLACEHOLDER_GOOGLE_USER_ID", // the "sub" field from Google — a numeric string like "1234567890"
		},

		JWTSecret:   "PLACEHOLDER_JWT_SECRET_REPLACE_WITH_32_PLUS_RANDOM_CHARS",
		JWTDuration: 24 * time.Hour,

		AppDeepLinkURI: "implodo://auth/complete",

		Port: "8080",
	}
}
