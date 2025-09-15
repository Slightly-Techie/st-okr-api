package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
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
	requestID := getRequestID(c)
	remoteIP := c.ClientIP()

	logger.Info("Company creation initiated",
		"request_id", requestID,
		"remote_ip", remoteIP,
	)

	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("Company creation failed - unauthorized",
			"request_id", requestID,
			"remote_ip", remoteIP,
			"reason", "no user_id in context",
		)
		response.Unauthorized(c, "Unauthorized request")
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Company creation failed - invalid request",
			"request_id", requestID,
			"user_id", userID.(string),
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	reqBody := dto.CreateCompanyRequest{
		Name:      body.Name,
		CreatorId: userID.(string),
	}

	logger.Info("Creating company",
		"request_id", requestID,
		"user_id", userID.(string),
		"company_name", body.Name,
	)

	data, err := ctrl.companyService.CreateCompany(reqBody)
	if err != nil {
		logger.Error("Company creation failed",
			"request_id", requestID,
			"user_id", userID.(string),
			"company_name", body.Name,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to create company", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Company created successfully",
		"request_id", requestID,
		"user_id", userID.(string),
		"company_id", data.ID,
		"company_name", data.Name,
	)

	response.Created(c, data, "Company created successfully")
}

func (ctrl *CompanyController) GetCompany(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Company retrieval requested",
		"request_id", requestID,
		"company_id", id,
		"user_id", userID,
	)

	data, err := ctrl.companyService.GetCompany("id", id)
	if err != nil {
		logger.Warn("Company not found",
			"request_id", requestID,
			"company_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Company not found")
		return
	}

	logger.Info("Company retrieved successfully",
		"request_id", requestID,
		"company_id", data.ID,
		"company_name", data.Name,
		"user_id", userID,
	)

	response.OK(c, data, "Company retrieved successfully")
}

func (ctrl *CompanyController) UpdateCompany(c *gin.Context) {
	var body dto.CreateCompanyRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	companyID := c.Param("id")

	logger.Info("Company update initiated",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Company update failed - invalid request",
			"request_id", requestID,
			"company_id", companyID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	logger.Info("Updating company",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
		"new_name", body.Name,
	)

	data, err := ctrl.companyService.UpdateCompany(body)
	if err != nil {
		logger.Error("Company update failed",
			"request_id", requestID,
			"company_id", companyID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update company", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Company updated successfully",
		"request_id", requestID,
		"company_id", data.ID,
		"company_name", data.Name,
		"user_id", userID,
	)

	response.OK(c, data, "Company updated successfully")
}

func (ctrl *CompanyController) DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Company deletion initiated",
		"request_id", requestID,
		"company_id", id,
		"user_id", userID,
	)

	err := ctrl.companyService.DeleteCompany(id)
	if err != nil {
		logger.Error("Company deletion failed",
			"request_id", requestID,
			"company_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to delete company", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Company deleted successfully",
		"request_id", requestID,
		"company_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Company deleted successfully")
}
