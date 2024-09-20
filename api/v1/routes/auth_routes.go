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

	authRoutes.GET("/:provider", authController.ContinueWithOAuth)         //localhost:8080/api/v1/google
	authRoutes.GET("/:provider/callback", authController.GetOAuthCallback) //localhost:8080/api/v1/google/callback
	authRoutes.GET("/logout/:provider", authController.LogoutWithOAuth)    //localhost:8080/api/v1/logout/google
}
