package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCompanyNotFound    = errors.New("no company exists with the provided credentials")
	ErrCompanyDBOperation = errors.New("database operation failed")
)

type CompanyRepository interface {
	GetDB() *gorm.DB
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

func (r *companyRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *companyRepository) GetByIdentifier(identifier, id string) (*models.Company, error) {
	var company models.Company

	res := r.db.Where(identifier, id).First(&company)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCompanyNotFound
		}
		log.Printf("error getting company by identifier: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrCompanyDBOperation, res.Error)
	}
	return &company, nil
}

func (r *companyRepository) Create(company *models.Company) (*models.Company, error) {
	res := r.db.Create(company)

	if res.Error != nil {
		log.Printf("error creating company: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrCompanyDBOperation, res.Error)
	}

	return company, nil
}

func (r *companyRepository) Update(company *models.Company) (*models.Company, error) {
	res := r.db.Save(company)

	if res.Error != nil {
		log.Printf("error updating company: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrCompanyDBOperation, res.Error)
	}

	return company, nil
}

func (r *companyRepository) Delete(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Company{})
	if res.Error != nil {
		log.Printf("error deleting company: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrCompanyDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no company found with id: %s", id)
		return ErrCompanyNotFound
	}
	return nil
}
