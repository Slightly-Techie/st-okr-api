package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/helper"
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyService interface {
	CreateCompany(r dto.CreateCompanyRequest) (*models.Company, error)
	GetCompany(ident, id string) (*models.Company, error)
	DeleteCompany(id string) error
	UpdateCompany(r dto.CreateCompanyRequest) (*models.Company, error)
}

type companyService struct {
	repo      repositories.CompanyRepository
	validator *validator.Validate
}

func NewCompanyService(repo repositories.CompanyRepository, validator *validator.Validate) CompanyService {
	return &companyService{
		repo:      repo,
		validator: validator,
	}
}

func (c *companyService) CreateCompany(r dto.CreateCompanyRequest) (*models.Company, error) {
	logger.Info("Company creation started",
		"company_name", r.Name,
		"creator_id", r.CreatorId,
	)

	if err := c.validator.Struct(r); err != nil {
		logger.Error("Company creation failed - validation error",
			"company_name", r.Name,
			"creator_id", r.CreatorId,
			"error", err.Error(),
		)
		return nil, err
	}

	var company models.Company

	err := c.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create company
		company = models.Company{
			ID:        uuid.NewString(),
			Name:      r.Name,
			Code:      helper.GenerateCompanyCode(r.Name, r.CreatorId),
			CreatorID: r.CreatorId,
		}

		logger.Debug("Creating company record",
			"company_id", company.ID,
			"company_name", company.Name,
			"company_code", company.Code,
		)

		if err := tx.Create(&company).Error; err != nil {
			logger.Error("Failed to create company record",
				"company_id", company.ID,
				"company_name", company.Name,
				"error", err.Error(),
			)
			return fmt.Errorf("failed to create company: %w", err)
		}

		// Create membership for the creator
		membership := models.Membership{
			ID:        uuid.NewString(),
			UserID:    r.CreatorId,
			CompanyID: company.ID,
			Role:      models.RoleAdmin, // Creator gets admin role
			Status:    models.StatusActive,
		}

		logger.Debug("Creating admin membership for creator",
			"membership_id", membership.ID,
			"user_id", membership.UserID,
			"company_id", membership.CompanyID,
			"role", membership.Role,
		)

		if err := tx.Create(&membership).Error; err != nil {
			logger.Error("Failed to create admin membership",
				"membership_id", membership.ID,
				"user_id", membership.UserID,
				"company_id", membership.CompanyID,
				"error", err.Error(),
			)
			return fmt.Errorf("failed to create membership: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Error("Company creation transaction failed",
			"company_name", r.Name,
			"creator_id", r.CreatorId,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	logger.Info("Company created successfully",
		"company_id", company.ID,
		"company_name", company.Name,
		"company_code", company.Code,
		"creator_id", company.CreatorID,
	)

	return &company, nil
}

func (c *companyService) GetCompany(ident, id string) (*models.Company, error) {
	logger.Debug("Retrieving company",
		"identifier_type", ident,
		"identifier_value", id,
	)

	company, err := c.repo.GetByIdentifier(ident, id)
	if err != nil {
		logger.Error("Failed to retrieve company",
			"identifier_type", ident,
			"identifier_value", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	logger.Info("Company retrieved successfully",
		"company_id", company.ID,
		"company_name", company.Name,
		"identifier_type", ident,
		"identifier_value", id,
	)

	return company, nil
}

func (c *companyService) UpdateCompany(r dto.CreateCompanyRequest) (*models.Company, error) {
	logger.Info("Company update started",
		"company_name", r.Name,
		"creator_id", r.CreatorId,
	)

	if err := c.validator.Struct(r); err != nil {
		logger.Error("Company update failed - validation error",
			"company_name", r.Name,
			"creator_id", r.CreatorId,
			"error", err.Error(),
		)
		return nil, err
	}

	company := models.Company{
		Name:      r.Name,
		CreatorID: r.CreatorId,
	}

	logger.Debug("Updating company record",
		"company_name", company.Name,
		"creator_id", company.CreatorID,
	)

	updatedCompany, err := c.repo.Update(&company)
	if err != nil {
		logger.Error("Failed to update company",
			"company_name", company.Name,
			"creator_id", company.CreatorID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	logger.Info("Company updated successfully",
		"company_id", updatedCompany.ID,
		"company_name", updatedCompany.Name,
		"creator_id", updatedCompany.CreatorID,
	)

	return updatedCompany, nil
}

func (c *companyService) DeleteCompany(id string) error {
	logger.Info("Company deletion started",
		"company_id", id,
	)

	if err := c.repo.Delete(id); err != nil {
		logger.Error("Failed to delete company",
			"company_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to delete company: %w", err)
	}

	logger.Info("Company deleted successfully",
		"company_id", id,
	)

	return nil
}
