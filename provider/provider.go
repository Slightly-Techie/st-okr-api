package provider

import (
	"github.com/Slightly-Techie/st-okr-api/internal/controllers"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Provider struct {
	UserController    *controllers.AuthController
	CompanyController *controllers.CompanyController
	DB                *gorm.DB
}

func NewProvider(db *gorm.DB, validator *validator.Validate) *Provider {
	// Initialize repository
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)

	// Initialize service
	userService := services.NewAuthService(userRepo, validator)
	companyService := services.NewCompanyService(companyRepo, validator)
	
	// Initialize handler
	userController := controllers.NewAuthController(userService)
	companyController := controllers.NewCompanyController(companyService)

	return &Provider{
		UserController:    userController,
		CompanyController: companyController,
		DB:                db,
	}
}
