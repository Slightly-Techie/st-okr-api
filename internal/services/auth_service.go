package services

import (
	"errors"
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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
	// Add provider to query parameters
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Get user data from OAuth provider
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to complete user auth: %w", err)
	}

	// Look for existing user
	var existingUser models.User
	result := s.repo.GetDB().Where("provider_id = ?", gothUser.UserID).First(&existingUser)

	// If user doesn't exist, create new user
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			existingUser = newUser
		} else {
			return nil, fmt.Errorf("database error: %w", result.Error)
	fmt.Println(existingUser)

	if existingUser == nil {
		data := models.User{
			ID:         uuid.NewString(),
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			AvatarURL:  user.AvatarURL,
			UserName:   user.NickName,
			ProviderID: user.UserID,
			Email:      user.Email,
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiry, err := auth.CreateJWTTokens(existingUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT tokens: %w", err)
	}

	// Create response
	response := &dto.AuthResponse{
	message.PublishMessage("sign_up", map[string]interface{}{
		"user_name": existingUser.UserName,
		"email":     existingUser.Email,
	})

	res := &dto.AuthResponse{
		FirstName:    existingUser.FirstName,
		LastName:     existingUser.LastName,
		UserName:     existingUser.UserName,
		Email:        existingUser.Email,
		ID:           existingUser.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiry,
	}

	return response, nil
}

func (s *authService) Logout(provider string, c *gin.Context) error {
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		return err
	}

	return nil
}
