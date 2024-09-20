package services

import (
	"github.com/Slightly-Techie/st-okr-api/api/v1/dto"
	"github.com/Slightly-Techie/st-okr-api/api/v1/models"
	"github.com/Slightly-Techie/st-okr-api/api/v1/repositories"
	"github.com/go-playground/validator/v10"
)

type CompanyService interface {
	CreateCompany(r dto.CreateCompanyRequest) (*models.Company, error)
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

// func (c *companyService) CreateCompany(r dto.CreateCompanyRequest) (*models.Company, error) {
// 	err := c.validator.Struct(r)
// 	if err != nil {
// 		return nil, err
// 	}


// }
