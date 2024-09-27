package services

import (
	"github.com/Slightly-Techie/st-okr-api/api/v1/dto"
	"github.com/Slightly-Techie/st-okr-api/api/v1/models"
	"github.com/Slightly-Techie/st-okr-api/api/v1/repositories"
	"github.com/Slightly-Techie/st-okr-api/internal/helper"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	err := c.validator.Struct(r)

	if err != nil {
		return nil, err
	}
	company := models.Company{
		ID:        uuid.NewString(),
		Name:      r.Name,
		Code:      helper.GenerateCompanyCode(r.Name, r.CreatorId),
		CreatorID: r.CreatorId,
	}
	resp, err := c.repo.Create(&company)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *companyService) GetCompany(ident, id string) (*models.Company, error) {
	// err := c.validator.Struct(r)
	// if err != nil {
	// 	return nil, err
	// }
	data, err := c.repo.GetByIdentifier(ident, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *companyService) UpdateCompany(r dto.CreateCompanyRequest) (*models.Company, error) {
	err := c.validator.Struct(r)
	if err != nil {
		return nil, err
	}
	company := models.Company{
		ID:        r.ID,
		Name:      r.Name,
		Code:      r.Code,
		CreatorID: r.CreatorId,
	}
	resp, err := c.repo.Update(&company)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *companyService) DeleteCompany(id string) error {
	// err := c.validator.Struct(r)
	// if err != nil {
	// 	return nil, err
	// }
	err := c.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
