package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrKeyResultNotFound    = errors.New("no key result exists with the provided details")
	ErrKeyResultDBOperation = errors.New("database operation failed")
)

type KeyResultRepository interface {
	GetDB() *gorm.DB
	Create(keyResult *models.KeyResult) (*models.KeyResult, error)
	GetByIdentifier(identifier, id string) (*models.KeyResult, error)
	ListByIdentifier(identifier, id string) ([]models.KeyResult, error)
	Update(keyResult *models.KeyResult) (*models.KeyResult, error)
	Delete(id string) error
}

type keyResultRepository struct {
	db *gorm.DB
}

func NewKeyResultRepository(db *gorm.DB) KeyResultRepository {
	return &keyResultRepository{db: db}
}

func (k *keyResultRepository) GetDB() *gorm.DB {
	return k.db
}

func (k *keyResultRepository) Create(keyResult *models.KeyResult) (*models.KeyResult, error) {
	res := k.db.Create(keyResult)
	if res.Error != nil {
		log.Printf("error creating team: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrKeyResultDBOperation, res.Error)
	}
	return keyResult, nil
}

func (k *keyResultRepository) GetByIdentifier(identifier, id string) (*models.KeyResult, error) {
	var keyResult models.KeyResult

	res := k.db.Where(identifier, id).First(&keyResult)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrKeyResultNotFound
		}
	}

	return &keyResult, nil
}

func (k *keyResultRepository) ListByIdentifier(identifier, id string) ([]models.KeyResult, error) {
	var keyResult []models.KeyResult

	res := k.db.Where(identifier, id).Find(&keyResult)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrKeyResultNotFound
		}
	}

	return keyResult, nil
}

func (k *keyResultRepository) Update(keyResult *models.KeyResult) (*models.KeyResult, error) {
	res := k.db.Save(keyResult)
	if res.Error != nil {
		log.Printf("error updating Key Result: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrKeyResultDBOperation, res.Error)
	}

	return keyResult, nil
}

func (k *keyResultRepository) Delete(id string) error {
	res := k.db.Where("id = ?", id).Delete(&models.Team{})
	if res.Error != nil {
		log.Printf("error deleting Key Result: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrKeyResultDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no Key Result found with id: %s", id)
		return ErrKeyResultNotFound
	}

	return nil
}
