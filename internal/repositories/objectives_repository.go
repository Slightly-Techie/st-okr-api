package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrObjectiveNotFound    = errors.New("no objective exists with the provided credentials")
	ErrObjectiveDBOperation = errors.New("database operation failed")
)

type ObjectivesRepository interface {
	GetDB() *gorm.DB
	GetByIdentifier(identifier, id string) (*models.Objectives, error)
	Create(objective *models.Objectives) (*models.Objectives, error)
	Update(objective *models.Objectives) (*models.Objectives, error)
	Delete(id string) error
}

type objectivesRepository struct {
	db *gorm.DB
}

func NewObjectivesRepository(db *gorm.DB) ObjectivesRepository {
	return &objectivesRepository{
		db: db,
	}
}

func (r *objectivesRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *objectivesRepository) GetByIdentifier(identifier, id string) (*models.Objectives, error) {
	var objective models.Objectives

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

func (r *objectivesRepository) Create(objective *models.Objectives) (*models.Objectives, error) {
	res := r.db.Create(objective)

	if res.Error != nil {
		log.Printf("error creating objective: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objective, nil
}

func (r *objectivesRepository) Update(objective *models.Objectives) (*models.Objectives, error) {
	res := r.db.Save(objective)

	if res.Error != nil {
		log.Printf("error updating objective: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}

	return objective, nil
}

func (r *objectivesRepository) Delete(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Objectives{})
	if res.Error != nil {
		log.Printf("error deleting objective: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrObjectiveDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no objectives found with id: %s", id)
		return ErrObjectiveNotFound
	}
	return nil
}
