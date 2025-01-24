package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ObjectivesService interface {
	CreateObjectives(r dto.CreateObjectivesRequest) (*models.Objectives, error)
	GetObjectives(ident, id string) (*models.Objectives, error)
	DeleteObjectives(id string) error
	UpdateObjectives(r dto.CreateObjectivesRequest) (*models.Objectives, error)
}

type objectivesService struct {
	repo      repositories.ObjectivesRepository
	validator *validator.Validate
}

func NewObjectivesService(repo repositories.ObjectivesRepository, validator *validator.Validate) ObjectivesService {
	return &objectivesService{
		repo:      repo,
		validator: validator,
	}
}

func (c *objectivesService) CreateObjectives(r dto.CreateObjectivesRequest) (*models.Objectives, error) {
	if err := c.validator.Struct(r); err != nil {
		return nil, err
	}

	var objective models.Objectives

	objective = models.Objectives{
		ID:        uuid.NewString(),
		Title:     r.Title,
		CreatorID: r.CreatorID,
		Deadline:  r.Deadline,
	}

	created, err := c.repo.Create(&objective)
	if err != nil {
		return nil, fmt.Errorf("failed to create objective: %w", err)
	}

	return created, nil

}

func (c *objectivesService) GetObjectives(ident, id string) (*models.Objectives, error) {
	objective, err := c.repo.GetByIdentifier(ident, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get objective: %w", err)
	}
	return objective, nil
}

func (c *objectivesService) UpdateObjectives(r dto.CreateObjectivesRequest) (*models.Objectives, error) {
	if err := c.validator.Struct(r); err != nil {
		return nil, err
	}

	objective := models.Objectives{
		Title:     r.Title,
		CreatorID: r.CreatorID,
	}

	updatedObjective, err := c.repo.Update(&objective)
	if err != nil {
		return nil, fmt.Errorf("failed to update objective: %w", err)
	}

	return updatedObjective, nil
}

func (c *objectivesService) DeleteObjectives(id string) error {
	if err := c.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete objective: %w", err)
	}
	return nil
}
