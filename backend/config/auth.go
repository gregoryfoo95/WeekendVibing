package config

import (
	"os"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthConfig struct {
	GoogleOAuth       *oauth2.Config
	JWTSecret         string
	JWTExpiration     time.Duration
	CookieDomain      string
	CookieSecure      bool
	CookieHttpOnly    bool
	CookieSameSite    string
}

func NewAuthConfig() *AuthConfig {
	googleOAuth := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-key-change-this-in-production"
	}

	jwtExpiration := 24 * time.Hour // 24 hours
	if exp := os.Getenv("JWT_EXPIRATION_HOURS"); exp != "" {
		if hours, err := time.ParseDuration(exp + "h"); err == nil {
			jwtExpiration = hours
		}
	}

	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = "localhost"
	}

	cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
	cookieHttpOnly := os.Getenv("COOKIE_HTTP_ONLY") != "false" // Default to true

	cookieSameSite := os.Getenv("COOKIE_SAME_SITE")
	if cookieSameSite == "" {
		cookieSameSite = "Lax"
	}

	return &AuthConfig{
		GoogleOAuth:       googleOAuth,
		JWTSecret:         jwtSecret,
		JWTExpiration:     jwtExpiration,
		CookieDomain:      cookieDomain,
		CookieSecure:      cookieSecure,
		CookieHttpOnly:    cookieHttpOnly,
		CookieSameSite:    cookieSameSite,
	}
} 