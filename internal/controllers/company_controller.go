package controllers

import (
	"fmt"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqBody := dto.CreateCompanyRequest{
		Name:      body.Name,
		CreatorId: userID.(string),
	}

	fmt.Println(reqBody)

	data, err := ctrl.companyService.CreateCompany(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company created successfully",
		"data":    data,
	})
}

func (ctrl *CompanyController) GetCompany(c *gin.Context) {
	// var body dto.CreateCompanyRequest

	id := c.Param("id")

	data, err := ctrl.companyService.GetCompany("id", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company found",
		"data":    data,
	})
}

func (ctrl *CompanyController) UpdateCompany(c *gin.Context) {
	var body dto.CreateCompanyRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.companyService.UpdateCompany(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company created successfully",
		"data":    data,
	})
}

func (ctrl *CompanyController) DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.companyService.DeleteCompany(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
