package routes

import (
	"github.com/Slightly-Techie/st-okr-api/api/v1/controllers"
	"github.com/Slightly-Techie/st-okr-api/api/v1/repositories"
	"github.com/Slightly-Techie/st-okr-api/api/v1/services"
	"github.com/Slightly-Techie/st-okr-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func AuthRoutes(r *gin.RouterGroup, validator *validator.Validate) {
	authRepo := repositories.NewUserRepository(database.DB)
	authService := services.NewAuthService(authRepo, validator)
	authController := controllers.NewAuthController(authService)
	authRoutes := r.Group("/auth")

	authRoutes.GET("/:provider", authController.ContinueWithOAuth)
	authRoutes.GET("/:provider/callback", authController.GetOAuthCallback)
	authRoutes.GET("/logout/:provider", authController.LogoutWithOAuth)
}
