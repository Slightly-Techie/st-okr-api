package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/helper"
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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
	if err := c.validator.Struct(r); err != nil {
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

		if err := tx.Create(&company).Error; err != nil {
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

		if err := tx.Create(&membership).Error; err != nil {
			return fmt.Errorf("failed to create membership: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &company, nil
}

func (c *companyService) GetCompany(ident, id string) (*models.Company, error) {
	company, err := c.repo.GetByIdentifier(ident, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}
	return company, nil
}

func (c *companyService) UpdateCompany(r dto.CreateCompanyRequest) (*models.Company, error) {
	if err := c.validator.Struct(r); err != nil {
		return nil, err
	}

	company := models.Company{
		Name:      r.Name,
		CreatorID: r.CreatorId,
	}

	updatedCompany, err := c.repo.Update(&company)
	if err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	return updatedCompany, nil
}

func (c *companyService) DeleteCompany(id string) error {
	if err := c.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
