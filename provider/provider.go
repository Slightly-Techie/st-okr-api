package provider

import (
	"github.com/Slightly-Techie/st-okr-api/internal/controllers"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Provider struct {
	UserController       *controllers.AuthController
	CompanyController    *controllers.CompanyController
	MembershipController *controllers.MembershipController
	TeamController       *controllers.TeamController
	DB                   *gorm.DB
}

func NewProvider(db *gorm.DB, validator *validator.Validate) *Provider {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	membershipRepo := repositories.NewMembershipRepository(db)
	teamRepo := repositories.NewTeamRepository(db)

	// Initialize services
	userService := services.NewAuthService(userRepo, validator)
	companyService := services.NewCompanyService(companyRepo, validator)
	membershipService := services.NewMembershipService(membershipRepo, validator)
	teamService := services.NewTeamService(teamRepo, validator)

	// Initialize controllers
	userController := controllers.NewAuthController(userService)
	companyController := controllers.NewCompanyController(companyService)
	membershipController := controllers.NewMembershipController(membershipService)
	teamController := controllers.NewTeamController(teamService)

	return &Provider{
		UserController:       userController,
		CompanyController:    companyController,
		MembershipController: membershipController,
		TeamController:       teamController,
		DB:                   db,
	}
}
