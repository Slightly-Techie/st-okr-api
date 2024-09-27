package routes

import (
	"github.com/Slightly-Techie/st-okr-api/api/v1/controllers"
	"github.com/Slightly-Techie/st-okr-api/api/v1/repositories"
	"github.com/Slightly-Techie/st-okr-api/api/v1/services"
	"github.com/Slightly-Techie/st-okr-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CompanyRoutes(r *gin.RouterGroup, validator *validator.Validate) {
	companyRepo := repositories.NewCompanyRepository(database.DB)
	companyService := services.NewCompanyService(companyRepo, validator)
	companyController := controllers.NewCompanyController(companyService)
	authRoutes := r.Group("/company")

	authRoutes.POST("/create", companyController.CreateCompany)       //localhost:8080/api/v1/company/create
	authRoutes.GET("/get/:id", companyController.GetCompany)          //localhost:8080/api/v1/company/get
	authRoutes.PUT("/update/:id", companyController.UpdateCompany)    //localhost:8080/api/v1/company/update
	authRoutes.DELETE("/delete/:id", companyController.DeleteCompany) //localhost:8080/api/v1/company/delete
}
