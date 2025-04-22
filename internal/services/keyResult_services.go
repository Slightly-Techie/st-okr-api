package services

import (
	"fmt"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/Slightly-Techie/st-okr-api/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type KeyResultService interface {
	CreateKeyResult(req dto.CreateKeyResultRequest) (*models.KeyResult, error)
	GetData(identifier, id string) (*models.KeyResult, error)
	UpdateKeyResult(req dto.UpdateKeyResultRequest) (*models.KeyResult, error)
	DeleteKeyResult(id string) error
	ListData(identifier, id string) ([]models.KeyResult, error)
}

type keyResultService struct {
	repo      repositories.KeyResultRepository
	validator *validator.Validate
}

func NewKeyResultService(repo repositories.KeyResultRepository, validator *validator.Validate) KeyResultService {
	validation.KeyResultValidators(validator)

	return &keyResultService{
		repo:      repo,
		validator: validator,
	}
}

func (k *keyResultService) CreateKeyResult(req dto.CreateKeyResultRequest) (*models.KeyResult, error) {
	if err := k.validator.Struct(req); err != nil {
		return nil, err
	}

	data := models.KeyResult{
		ID:           uuid.NewString(),
		ObjectiveID:  req.ObjectiveID,
		Title:        req.Title,
		Description:  req.Description,
		MetricType:   models.MetricType(req.MetricType),
		CurrentValue: req.CurrentValue,
		TargetValue:  req.TargetValue,
		AssigneeType: models.AssigneeType(req.AssigneeType),
		AssigneeID:   req.AssigneeID,
		StartDate:    req.StartDate,
		DueDate:      req.DueDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := validation.ValidateMetricValues(&data); err != nil {
		return nil, err
	}

	data.UpdateProgress()
	data.UpdateStatus()

	created, err := k.repo.Create(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Result: %w", err)
	}

	return created, nil
}

func (k *keyResultService) GetData(identifier, id string) (*models.KeyResult, error) {
	res, err := k.repo.GetByIdentifier(identifier, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %v", err)
	}

	return res, nil
}

func (k *keyResultService) UpdateKeyResult(req dto.UpdateKeyResultRequest) (*models.KeyResult, error) {
	if err := k.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	existing, err := k.repo.GetByIdentifier("id", req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find key result: %w", err)
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.CurrentValue = req.CurrentValue
	existing.TargetValue = req.TargetValue
	existing.MetricType = req.MetricType
	existing.AssigneeType = req.AssigneeType
	existing.AssigneeID = req.AssigneeID
	existing.StartDate = req.StartDate
	existing.DueDate = req.DueDate

	if err := validation.ValidateMetricValues(existing); err != nil {
		return nil, err
	}

	existing.UpdateProgress()
	existing.UpdateStatus()

	updatedData, err := k.repo.Update(existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update key Result: %v", err)
	}

	return updatedData, nil
}

func (k *keyResultService) DeleteKeyResult(id string) error {
	if err := k.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete key Result: %v", err)
	}
	return nil
}

func (k *keyResultService) ListData(identifier, objId string) ([]models.KeyResult, error) {
	keys, err := k.repo.ListByIdentifier(identifier, objId)
	if err != nil {
		return nil, fmt.Errorf("failed to list data: %v", err)
	}

	return keys, nil
}

// func (k *keyResultService) ListAssigneeKeyResults(identifier, userId string) (*models.KeyResult, error) {
// 	assignee, err := k.repo.GetByIdentifier(identifier, userId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get Assignee's Key Results: %v", err)
// 	}

// 	return assignee, nil
// }
