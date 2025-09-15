package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

// Helper functions for logging
func getRequestID(c *gin.Context) string {
	if reqID, exists := c.Get("request_id"); exists {
		return reqID.(string)
	}
	return "unknown"
}

func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

func maskEmail(email string) string {
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

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) ContinueWithOAuth(c *gin.Context) {
	provider := c.Param("provider")
	requestID := getRequestID(c)
	remoteIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	logger.Info("OAuth authentication initiated",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
		"user_agent", userAgent,
	)

	if provider == "" {
		logger.Error("OAuth authentication failed - missing provider",
			"request_id", requestID,
			"remote_ip", remoteIP,
			"error", "provider parameter is empty",
		)
		response.BadRequest(c, "OAuth provider not specified", map[string]string{
			"provider": "Provider parameter is required",
		})
		return
	}

	logger.Debug("Delegating to auth service",
		"request_id", requestID,
		"provider", provider,
	)

	ctrl.authService.AuthHandler(provider, c)
}

func (ctrl *AuthController) GetOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	requestID := getRequestID(c)
	remoteIP := c.ClientIP()

	logger.Info("OAuth callback received",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
	)

	authRes, err := ctrl.authService.GetAuthCallback(provider, c)
	if err != nil {
		logger.Error("OAuth authentication failed during callback",
			"request_id", requestID,
			"provider", provider,
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		response.InternalError(c, "Authentication failed")
		return
	}

	// Extract user info for logging (safely)
	userID := ""
	userEmail := ""
	if authRes != nil {
		userID = authRes.ID
		userEmail = authRes.Email
	}

	logger.Info("OAuth authentication successful",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
		"user_id", userID,
		"user_email", maskEmail(userEmail),
	)

	userData, err := json.Marshal(authRes)
	if err != nil {
		logger.Error("Failed to marshal user data",
			"request_id", requestID,
			"provider", provider,
			"user_id", userID,
			"error", err.Error(),
		)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/error?message=failed to process user data", config.ENV.DBHost))
		return
	}

	logger.Debug("Setting authentication cookie",
		"request_id", requestID,
		"user_id", userID,
		"domain", config.ENV.DBHost,
	)

	c.SetCookie(
		"user_data",       // cookie name
		string(userData),  // cookie value
		3600*24,           // expiration time (24 hours)
		"/",               // path
		config.ENV.DBHost, // domain
		true,              // secure
		true,              // httpOnly
	)

	logger.Info("User redirected to dashboard",
		"request_id", requestID,
		"user_id", userID,
		"redirect_url", "http://localhost:5173/dashboard",
	)

	// Replace the url with the actual url
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/dashboard")
}

func (ctrl *AuthController) LogoutWithOAuth(c *gin.Context) {
	provider := c.Param("provider")
	requestID := getRequestID(c)
	remoteIP := c.ClientIP()
	userID := getUserID(c)

	logger.Info("User logout initiated",
		"request_id", requestID,
		"provider", provider,
		"remote_ip", remoteIP,
		"user_id", userID,
	)

	err := ctrl.authService.Logout(provider, c)
	if err != nil {
		logger.Error("Logout failed",
			"request_id", requestID,
			"provider", provider,
			"user_id", userID,
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		response.InternalError(c, "Logout failed")
		return
	}

	logger.Info("User logout successful",
		"request_id", requestID,
		"provider", provider,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
