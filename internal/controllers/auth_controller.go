package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) ContinueWithOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"error": "Provider not specified"})
		return
	}

	ctrl.authService.AuthHandler(provider, c)
}

func (ctrl *AuthController) GetOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")

	authRes, err := ctrl.authService.GetAuthCallback(provider, c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userData, err := json.Marshal(authRes)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/error?message=failed to process user data", config.ENV.DBHost))
		return
	}
	c.SetCookie(
		"user_data",       // cookie name
		string(userData),  // cookie value
		3600*24,           // expiration time (24 hours)
		"/",               // path
		config.ENV.DBHost, // domain
		true,              // secure
		true,              // httpOnly
	)

	// Replace the url with the actual url
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/dashboard")
}

func (ctrl *AuthController) LogoutWithOAuth(c *gin.Context) {
	provider := c.Param("provider")

	err := ctrl.authService.Logout(provider, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
