package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"fithero-backend/config"
	"fithero-backend/services"
	"fithero-backend/middleware"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
	authConfig  *config.AuthConfig
}

func NewAuthController(authService *services.AuthService, authConfig *config.AuthConfig) *AuthController {
	return &AuthController{
		authService: authService,
		authConfig:  authConfig,
	}
}

// GoogleLogin initiates Google OAuth flow
func (ac *AuthController) GoogleLogin(c *gin.Context) {
	// Generate state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate state",
		})
		return
	}

	// Store state in session/cookie for verification
	c.SetCookie(
		"oauth_state",
		state,
		300, // 5 minutes
		"/",
		ac.authConfig.CookieDomain,
		ac.authConfig.CookieSecure,
		ac.authConfig.CookieHttpOnly,
	)

	// Get Google OAuth URL
	url := ac.authService.GetGoogleAuthURL(state)

	// Return the OAuth URL as JSON for frontend to handle redirect
	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
	})
}

// GoogleCallback handles Google OAuth callback
func (ac *AuthController) GoogleCallback(c *gin.Context) {
	// Verify state parameter
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter",
		})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", ac.authConfig.CookieDomain, ac.authConfig.CookieSecure, ac.authConfig.CookieHttpOnly)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code required",
		})
		return
	}

	// Handle the callback
	authResponse, err := ac.authService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Authentication failed",
			"message": err.Error(),
		})
		return
	}

	// Set secure HTTP-only cookie with JWT token
	ac.setAuthCookie(c, authResponse.Token)

	// Redirect to frontend auth callback to complete the flow
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/auth-callback")
}

// Logout clears the authentication cookie
func (ac *AuthController) Logout(c *gin.Context) {
	// Clear the auth cookie
	c.SetCookie(
		"auth_token",
		"",
		-1, // Expire immediately
		"/",
		ac.authConfig.CookieDomain,
		ac.authConfig.CookieSecure,
		ac.authConfig.CookieHttpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// RefreshToken generates a new JWT token for the authenticated user
func (ac *AuthController) RefreshToken(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	newToken, err := ac.authService.RefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh token",
		})
		return
	}

	// Set the new token in cookie
	ac.setAuthCookie(c, newToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
	})
}

// Me returns the current authenticated user's profile
func (ac *AuthController) Me(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// CheckAuth verifies if user is authenticated
func (ac *AuthController) CheckAuth(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user":          user,
	})
}

// Private helper methods

func (ac *AuthController) setAuthCookie(c *gin.Context, token string) {
	sameSite := http.SameSiteLaxMode
	switch ac.authConfig.CookieSameSite {
	case "Strict":
		sameSite = http.SameSiteStrictMode
	case "None":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		"auth_token",
		token,
		int(ac.authConfig.JWTExpiration.Seconds()),
		"/",
		ac.authConfig.CookieDomain,
		ac.authConfig.CookieSecure,
		ac.authConfig.CookieHttpOnly,
	)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
} 