package services

import (
	"errors"
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/message"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	auth "github.com/Slightly-Techie/st-okr-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthService interface {
	AuthHandler(provider string, c *gin.Context)
	GetAuthCallback(provider string, c *gin.Context) (*dto.AuthResponse, error)
	Logout(provider string, c *gin.Context) error
}

type authService struct {
	repo      repositories.UserRepository
	validator *validator.Validate
}

func NewAuthService(repo repositories.UserRepository, validator *validator.Validate) AuthService {
	return &authService{
		repo:      repo,
		validator: validator,
	}
}

func (s *authService) AuthHandler(provider string, c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (s *authService) GetAuthCallback(provider string, c *gin.Context) (*dto.AuthResponse, error) {
	requestID := getRequestIDFromContext(c)
	remoteIP := c.ClientIP()

	logger.Info("Processing OAuth callback",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
	)

	// Add provider to query parameters
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Get user data from OAuth provider
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		logger.Error("Failed to complete OAuth user authentication",
			"request_id", requestID,
			"provider", provider,
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to complete user auth: %w", err)
	}

	logger.Debug("OAuth user data received",
		"request_id", requestID,
		"provider", provider,
		"provider_user_id", gothUser.UserID,
		"email", maskEmailForLog(gothUser.Email),
		"nickname", gothUser.NickName,
	)

	// Look for existing user
	logger.Debug("Looking up existing user",
		"request_id", requestID,
		"provider_user_id", gothUser.UserID,
	)

	var existingUser models.User
	result := s.repo.GetDB().Where("provider_id = ?", gothUser.UserID).First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Info("Creating new user account",
				"request_id", requestID,
				"provider", provider,
				"provider_user_id", gothUser.UserID,
				"email", maskEmailForLog(gothUser.Email),
				"username", gothUser.NickName,
			)

			// If user doesn't exist, create new user
			newUser := models.User{
				ID:         uuid.NewString(),
				FirstName:  gothUser.FirstName,
				LastName:   gothUser.LastName,
				AvatarURL:  gothUser.AvatarURL,
				UserName:   gothUser.NickName,
				ProviderID: gothUser.UserID,
				Email:      gothUser.Email,
			}

			if err := s.repo.GetDB().Create(&newUser).Error; err != nil {
				logger.Error("Failed to create new user",
					"request_id", requestID,
					"provider_user_id", gothUser.UserID,
					"email", maskEmailForLog(gothUser.Email),
					"error", err.Error(),
				)
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			logger.Info("New user created successfully",
				"request_id", requestID,
				"user_id", newUser.ID,
				"provider", provider,
				"email", maskEmailForLog(newUser.Email),
			)

			existingUser = newUser
		} else {
			logger.Error("Database error during user lookup",
				"request_id", requestID,
				"provider_user_id", gothUser.UserID,
				"error", result.Error.Error(),
			)
			return nil, fmt.Errorf("database error: %w", result.Error)
		}
	} else {
		logger.Info("Existing user found",
			"request_id", requestID,
			"user_id", existingUser.ID,
			"provider", provider,
			"email", maskEmailForLog(existingUser.Email),
		)
	}

	// Generate JWT tokens
	logger.Debug("Generating JWT tokens",
		"request_id", requestID,
		"user_id", existingUser.ID,
	)

	accessToken, refreshToken, expiry, err := auth.CreateJWTTokens(existingUser.ID)
	if err != nil {
		logger.Error("Failed to create JWT tokens",
			"request_id", requestID,
			"user_id", existingUser.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create JWT tokens: %w", err)
	}

	logger.Info("JWT tokens generated successfully",
		"request_id", requestID,
		"user_id", existingUser.ID,
		"expires_in", expiry,
	)

	// Publish sign-up message
	logger.Debug("Publishing sign-up message",
		"request_id", requestID,
		"user_id", existingUser.ID,
	)

	message.PublishMessage("sign_up", map[string]any{
		"user_name": existingUser.UserName,
		"email":     existingUser.Email,
	})

	// Create response
	response := &dto.AuthResponse{
		FirstName:    existingUser.FirstName,
		LastName:     existingUser.LastName,
		UserName:     existingUser.UserName,
		Email:        existingUser.Email,
		ID:           existingUser.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiry,
	}

	logger.Info("OAuth authentication completed successfully",
		"request_id", requestID,
		"user_id", existingUser.ID,
		"provider", provider,
		"email", maskEmailForLog(existingUser.Email),
	)

	return response, nil
}

func (s *authService) Logout(provider string, c *gin.Context) error {
	requestID := getRequestIDFromContext(c)
	remoteIP := c.ClientIP()

	logger.Info("Processing logout request",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
	)

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		logger.Error("Failed to logout user",
			"request_id", requestID,
			"provider", provider,
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		return err
	}

	logger.Info("User logout completed successfully",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
	)

	return nil
}

// Helper functions
func getRequestIDFromContext(c *gin.Context) string {
	if reqID, exists := c.Get("request_id"); exists {
		return reqID.(string)
	}
	return "unknown"
}

func maskEmailForLog(email string) string {
	if email == "" {
		return ""
	}
	if len(email) <= 3 {
		return "***"
	}
	// Show first 2 characters and domain
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}
	if atIndex == -1 {
		return "***"
	}
	if atIndex < 2 {
		return "**" + email[atIndex:]
	}
	return email[:2] + "***" + email[atIndex:]
}
