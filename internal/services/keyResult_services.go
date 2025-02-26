package services

import (
	"fmt"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type KeyResultService interface {
	CreateKeyResult(req dto.CreateKeyResultRequest) (*models.KeyResult, error)
	GetData(identifier, id string) (*models.KeyResult, error)
	UpdateKeyResult(req dto.UpdateKeyResultRequest) (*models.KeyResult, error)
	DeleteKeyResult(id string) error
	ListData(identifier, id string) (*models.KeyResult, error)
}

type keyResultService struct {
	repo      repositories.KeyResultRepository
	validator *validator.Validate
}

func NewKeyResultService(repo repositories.KeyResultRepository, validator *validator.Validate) KeyResultService {
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
		MetricType:   req.MetricType,
		CurrentValue: req.CurrentValue,
		TargetValue:  req.TargetValue,
		AssigneeType: req.AssigneeType,
		AssigneeID:   req.AssigneeID,
		StartDate:    req.StartDate,
		DueDate:      req.DueDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
		return nil, err
	}

	data := models.KeyResult{
		Title:        req.Title,
		Description:  req.Description,
		CurrentValue: req.CurrentValue,
		// Progress: req.
		Status:       req.Status,
		AssigneeType: req.AssigneeType,
		AssigneeID:   req.AssigneeID,
		DueDate:      req.DueDate,
		UpdatedAt:    time.Now(),
	}
	data.UpdateProgress()
	data.UpdateStatus()

	updatedData, err := k.repo.Update(&data)
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

func (k *keyResultService) ListData(identifier, objId string) (*models.KeyResult, error) {
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
