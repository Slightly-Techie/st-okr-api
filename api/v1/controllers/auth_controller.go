package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Slightly-Techie/st-okr-api/api/v1/services"
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

	user, err := ctrl.authService.GetAuthCallback(provider, c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"user": user})
	c.Redirect(http.StatusFound, "http://localhost:5174")
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
