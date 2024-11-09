package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	GetByIdentifier(identifier, id string) (*models.Company, error)
	Create(company *models.Company) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id string) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

func (r *companyRepository) GetByIdentifier(identifier, id string) (*models.Company, error) {
	var company models.Company

	res := r.db.Where(identifier, id).First(&company)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no company exists with the provided credentials")
		}
		log.Println("error getting company by identifier: ", res.Error)
		return nil, res.Error
	}
	return &company, nil
}

func (r *companyRepository) Create(company *models.Company) (*models.Company, error) {
	res := r.db.Create(&company)

	if res.Error != nil {
		log.Println("error creating company: ", res.Error)
		return nil, res.Error
	}

	return company, nil
}

func (r *companyRepository) Update(company *models.Company) (*models.Company, error) {
	res := r.db.Save(&company)

	if res.Error != nil {
		log.Println("error updating company: ", res.Error)
		return nil, res.Error
	}

	return company, nil
}

func (r *companyRepository) Delete(id string) error {
	var company models.Company

	res := r.db.Where("id = ?", id).Delete(&company)
	if res.Error != nil {
		log.Println("error deleting company: ", res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Println("no company found with the given id")
		return fmt.Errorf("no company found with id: %s", id)
	}
	return nil
}