package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/api/v1/dto"
	"github.com/Slightly-Techie/st-okr-api/api/v1/services"
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
	var body dto.CreateCompanyRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.companyService.CreateCompany(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, _ := json.Marshal(data)
	c.JSON(http.StatusOK, gin.H{
		"message": "Company created successfully",
		"data":    string(resp),
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

	resp, _ := json.Marshal(data)
	c.JSON(http.StatusOK, gin.H{
		"message": "Company found",
		"data":    string(resp),
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
	resp, _ := json.Marshal(data)
	c.JSON(http.StatusOK, gin.H{
		"message": "Company created successfully",
		"data":    string(resp),
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
