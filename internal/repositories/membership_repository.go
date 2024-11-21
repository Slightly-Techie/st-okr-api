package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrMembershipNotFound    = errors.New("no membership exists with the provided credentials")
	ErrMembershipDBOperation = errors.New("database operation failed")
)

type MembershipRepository interface {
	GetDB() *gorm.DB
	GetByIdentifier(identifier, id string) (*models.Membership, error)
	Create(membership *models.Membership) (*models.Membership, error)
	Update(membership *models.Membership) (*models.Membership, error)
	Delete(id string) error
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{
		db: db,
	}
}

func (r *membershipRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *membershipRepository) GetByIdentifier(identifier, id string) (*models.Membership, error) {
	var membership models.Membership

	res := r.db.Where(identifier, id).First(&membership)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrMembershipNotFound
		}
		log.Printf("error getting membership by identifier: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrMembershipDBOperation, res.Error)
	}
	return &membership, nil
}

func (r *membershipRepository) Create(membership *models.Membership) (*models.Membership, error) {
	res := r.db.Create(membership)

	if res.Error != nil {
		log.Printf("error creating membership: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrMembershipDBOperation, res.Error)
	}

	return membership, nil
}

func (r *membershipRepository) Update(membership *models.Membership) (*models.Membership, error) {
	res := r.db.Save(membership)

	if res.Error != nil {
		log.Printf("error updating membership: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrMembershipDBOperation, res.Error)
	}

	return membership, nil
}

func (r *membershipRepository) Delete(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Membership{})
	if res.Error != nil {
		log.Printf("error deleting membership: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrMembershipDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no membership found with id: %s", id)
		return ErrMembershipNotFound
	}
	return nil
}