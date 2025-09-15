package controllers

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type CompanyController struct {
	companyService services.CompanyService
}

func NewCompanyController(companyService services.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

func (ctrl *CompanyController) CreateCompany(c *gin.Context) {
	type CCBody struct {
		Name string `json:"name"`
	}
	var body CCBody

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized request")
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	reqBody := dto.CreateCompanyRequest{
		Name:      body.Name,
		CreatorId: userID.(string),
	}

	fmt.Println(reqBody)

	data, err := ctrl.companyService.CreateCompany(reqBody)
	if err != nil {
		response.BadRequest(c, "Failed to create company", map[string]string{
			"service": err.Error(),
		})
		return
	}
	response.Created(c, data, "Company created successfully")
}

func (ctrl *CompanyController) GetCompany(c *gin.Context) {
	id := c.Param("id")

	data, err := ctrl.companyService.GetCompany("id", id)
	if err != nil {
		response.NotFound(c, "Company not found")
		return
	}

	response.OK(c, data, "Company retrieved successfully")
}

func (ctrl *CompanyController) UpdateCompany(c *gin.Context) {
	var body dto.CreateCompanyRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	data, err := ctrl.companyService.UpdateCompany(body)
	if err != nil {
		response.BadRequest(c, "Failed to update company", map[string]string{
			"service": err.Error(),
		})
		return
	}
	response.OK(c, data, "Company updated successfully")
}

func (ctrl *CompanyController) DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.companyService.DeleteCompany(id)
	if err != nil {
		response.BadRequest(c, "Failed to delete company", map[string]string{
			"service": err.Error(),
		})
		return
	}
	response.OK(c, nil, "Company deleted successfully")
}
