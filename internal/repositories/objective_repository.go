package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrObjectiveNotFound    = errors.New("no objective exists with the provided details")
	ErrObjectiveDBOperation = errors.New("database operation failed")
)

type ObjectiveRepository interface {
	GetDB() *gorm.DB
	Create(objective *models.Objective) (*models.Objective, error)
	GetByIdentifier(identifier, id string) (*models.Objective, error)
	GetWithKeyResults(id string) (*models.Objective, error)
	ListByIdentifier(identifier, id string) ([]models.Objective, error)
	ListByCompany(companyID string) ([]models.Objective, error)
	ListByTeam(teamID string) ([]models.Objective, error)
	ListByOwner(ownerID string) ([]models.Objective, error)
	Update(objective *models.Objective) (*models.Objective, error)
	Delete(id string) error
}

type objectiveRepository struct {
	db *gorm.DB
}

func NewObjectiveRepository(db *gorm.DB) ObjectiveRepository {
	return &objectiveRepository{db: db}
}

func (r *objectiveRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *objectiveRepository) Create(objective *models.Objective) (*models.Objective, error) {
	res := r.db.Create(objective)
	if res.Error != nil {
		log.Printf("error creating objective: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}
	return objective, nil
}

func (r *objectiveRepository) GetByIdentifier(identifier, id string) (*models.Objective, error) {
	var objective models.Objective

	res := r.db.Where(identifier, id).First(&objective)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrObjectiveNotFound
		}
		log.Printf("error getting objective by identifier: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return &objective, nil
}

func (r *objectiveRepository) GetWithKeyResults(id string) (*models.Objective, error) {
	var objective models.Objective

	res := r.db.Preload("KeyResults").Where("id = ?", id).First(&objective)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrObjectiveNotFound
		}
		log.Printf("error getting objective with key results: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return &objective, nil
}

func (r *objectiveRepository) ListByIdentifier(identifier, id string) ([]models.Objective, error) {
	var objectives []models.Objective

	res := r.db.Where(identifier, id).Find(&objectives)
	if res.Error != nil {
		log.Printf("error listing objectives by identifier: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objectives, nil
}

func (r *objectiveRepository) ListByCompany(companyID string) ([]models.Objective, error) {
	var objectives []models.Objective

	res := r.db.Where("company_id = ?", companyID).Find(&objectives)
	if res.Error != nil {
		log.Printf("error listing objectives by company: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objectives, nil
}

func (r *objectiveRepository) ListByTeam(teamID string) ([]models.Objective, error) {
	var objectives []models.Objective

	res := r.db.Where("team_id = ?", teamID).Find(&objectives)
	if res.Error != nil {
		log.Printf("error listing objectives by team: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objectives, nil
}

func (r *objectiveRepository) ListByOwner(ownerID string) ([]models.Objective, error) {
	var objectives []models.Objective

	res := r.db.Where("owner_id = ?", ownerID).Find(&objectives)
	if res.Error != nil {
		log.Printf("error listing objectives by owner: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objectives, nil
}

func (r *objectiveRepository) Update(objective *models.Objective) (*models.Objective, error) {
	res := r.db.Save(objective)
	if res.Error != nil {
		log.Printf("error updating objective: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objective, nil
}

func (r *objectiveRepository) Delete(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Objective{})
	if res.Error != nil {
		log.Printf("error deleting objective: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no objective found with id: %s", id)
		return ErrObjectiveNotFound
	}

	return nil
}