package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"fithero-backend/services"
	"fithero-backend/models"

	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "user"
	UserIDContextKey = "user_id"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from Authorization header
		token := extractTokenFromHeader(c)
		
		// If not in header, try to get from cookie
		if token == "" {
			if cookie, err := c.Request.Cookie("auth_token"); err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user details (optional - for additional validation)
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User account is disabled",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, claims.UserID)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require authentication
func OptionalAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from Authorization header
		token := extractTokenFromHeader(c)
		
		// If not in header, try to get from cookie
		if token == "" {
			if cookie, err := c.Request.Cookie("auth_token"); err == nil {
				token = cookie.Value
			}
		}

		if token != "" {
			// Validate the token
			if claims, err := authService.ValidateJWT(token); err == nil {
				// Get user details
				if user, err := authService.GetUserByID(claims.UserID); err == nil && user.IsActive {
					// Set user context
					c.Set(UserContextKey, user)
					c.Set(UserIDContextKey, claims.UserID)
				}
			}
		}

		c.Next()
	}
}

// extractTokenFromHeader extracts JWT token from Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Authorization header format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetCurrentUser returns the current authenticated user from context
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	
	userModel, ok := user.(*models.User)
	return userModel, ok
}

// GetCurrentUserID returns the current authenticated user ID from context
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return 0, false
	}
	
	id, ok := userID.(uint)
	return id, ok
}

// RequireOwnership middleware ensures the user can only access their own resources
func RequireOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Get the resource user ID from URL parameter
		resourceUserIDStr := c.Param("user_id")
		if resourceUserIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID required in URL",
			})
			c.Abort()
			return
		}

		// Convert to uint (simple implementation - you might want better parsing)
		var resourceUserID uint
		if _, err := fmt.Sscanf(resourceUserIDStr, "%d", &resourceUserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Check if the current user owns the resource
		if currentUserID != resourceUserID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: You can only access your own resources",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 