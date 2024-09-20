package services

import (
	"github.com/Slightly-Techie/st-okr-api/api/v1/dto"
	"github.com/Slightly-Techie/st-okr-api/api/v1/models"
	"github.com/Slightly-Techie/st-okr-api/api/v1/repositories"
	"github.com/Slightly-Techie/st-okr-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
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
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		return nil, err
	}

	

	existingUser, err := s.repo.GetByIdentifier("provider_id", user.UserID)

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

		usr, err := s.repo.Create(&data)
		if err != nil {
			return nil, err
		}

		existingUser = usr
	}

	access_token, refresh_token, expiry, err := auth.CreateJWTTokens(existingUser.ID)
	if err != nil {
		return nil, err
	}

	res := &dto.AuthResponse{
		FirstName:    existingUser.FirstName,
		LastName:     existingUser.LastName,
		UserName:     existingUser.UserName,
		Email:        existingUser.Email,
		ID:           existingUser.ID,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiresIn:    expiry,
	}
	return res, nil
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
