package repositories

import (
	"errors"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByIdentifier(identifier, id string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByIdentifier(identifier, id string) (*models.User, error) {
	var user models.User

	res := r.db.Where(identifier, id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no user exists with the provided credentials")
		}
		log.Println("error getting user by identifier: ", res.Error)
		return nil, res.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	res := r.db.Create(&user)

	if res.Error != nil {
		log.Println("error creating user: ", res.Error)
		return nil, res.Error
	}

	return user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	res := r.db.Save(&user)

	if res.Error != nil {
		log.Println("error updating user: ", res.Error)
		return nil, res.Error
	}

	return user, nil
}

func (r *userRepository) Delete(id string) error {
	var user models.User

	res := r.db.Where("id = ?", id).Delete(&user)
	if res.Error != nil {
		log.Println("error deleting user: ", res.Error)
		return res.Error
	}

	return nil
}