package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fithero-backend/config"
	"fithero-backend/models"
	"fithero-backend/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthService struct {
	userRepo   repositories.UserRepositoryInterface
	authConfig *config.AuthConfig
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, authConfig *config.AuthConfig) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		authConfig: authConfig,
	}
}

// GetGoogleAuthURL generates the Google OAuth URL for authentication
func (s *AuthService) GetGoogleAuthURL(state string) string {
	return s.authConfig.GoogleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleGoogleCallback processes the Google OAuth callback
func (s *AuthService) HandleGoogleCallback(code string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.authConfig.GoogleOAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	googleUser, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}

	// Find or create user
	user, err := s.findOrCreateUser(googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user.ID, &models.UpdateUserRequest{}); err != nil {
		// Don't fail auth if we can't update last login time
		fmt.Printf("Warning: failed to update last login time for user %d: %v\n", user.ID, err)
	}

	// Generate JWT token
	jwtToken, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return &models.AuthResponse{
		User:  *user,
		Token: jwtToken,
	}, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (s *AuthService) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserByID retrieves a user by ID (used for authorization checks)
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// RefreshToken generates a new JWT token for the user
func (s *AuthService) RefreshToken(userID uint) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return "", fmt.Errorf("user account is disabled")
	}

	return s.generateJWT(user)
}

// Private methods

func (s *AuthService) getGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser models.GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (s *AuthService) findOrCreateUser(googleUser *models.GoogleUserInfo) (*models.User, error) {
	// Try to find existing user by Google ID
	if existingUser, err := s.userRepo.GetByGoogleID(googleUser.ID); err == nil && existingUser != nil {
		// Update user info from Google (in case they changed their profile)
		updateReq := &models.UpdateUserRequest{
			Email:     &googleUser.Email,
			FirstName: &googleUser.GivenName,
			LastName:  &googleUser.FamilyName,
		}
		if err := s.userRepo.Update(existingUser.ID, updateReq); err != nil {
			fmt.Printf("Warning: failed to update user info for user %d: %v\n", existingUser.ID, err)
		}
		return existingUser, nil
	}

	// Try to find existing user by email
	if existingUser, err := s.userRepo.GetByEmail(googleUser.Email); err == nil && existingUser != nil {
		// Link the Google account to existing user
		existingUser.GoogleID = googleUser.ID
		existingUser.FirstName = googleUser.GivenName
		existingUser.LastName = googleUser.FamilyName
		existingUser.Picture = googleUser.Picture
		
		if err := s.userRepo.Update(existingUser.ID, &models.UpdateUserRequest{
			FirstName: &googleUser.GivenName,
			LastName:  &googleUser.FamilyName,
		}); err != nil {
			fmt.Printf("Warning: failed to link Google account for user %d: %v\n", existingUser.ID, err)
		}
		
		return existingUser, nil
	}

	// Create new user
	newUser := &models.User{
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Username:  generateUsernameFromEmail(googleUser.Email),
		FirstName: googleUser.GivenName,
		LastName:  googleUser.FamilyName,
		Picture:   googleUser.Picture,
		Level:     1,
		Points:    0,
		Character: "Rookie Hero",
		JobTitle:  "Fitness Novice",
		IsActive:  true,
	}

	return s.userRepo.Create(newUser)
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := &models.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authConfig.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fithero-backend",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.authConfig.JWTSecret))
}

// generateUsernameFromEmail creates a username from email address
func generateUsernameFromEmail(email string) string {
	// Simple implementation - take part before @ and add random suffix if needed
	// In production, you might want to check for uniqueness and add numbers
	atIndex := 0
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}
	if atIndex > 0 {
		return email[:atIndex]
	}
	return "user"
} 