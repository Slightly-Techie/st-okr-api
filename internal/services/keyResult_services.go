package services

import (
	"fmt"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
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
	logger.Info("Key result creation started",
		"title", req.Title,
		"objective_id", req.ObjectiveID,
		"assignee_id", req.AssigneeID,
		"assignee_type", req.AssigneeType,
		"metric_type", req.MetricType,
	)

	if err := k.validator.Struct(req); err != nil {
		logger.Error("Key result creation failed - validation error",
			"title", req.Title,
			"objective_id", req.ObjectiveID,
			"error", err.Error(),
		)
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

	logger.Debug("Validating metric values",
		"key_result_id", data.ID,
		"metric_type", data.MetricType,
		"current_value", data.CurrentValue,
		"target_value", data.TargetValue,
	)

	if err := validation.ValidateMetricValues(&data); err != nil {
		logger.Error("Key result creation failed - metric validation error",
			"key_result_id", data.ID,
			"title", data.Title,
			"metric_type", data.MetricType,
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug("Calculating progress and status",
		"key_result_id", data.ID,
		"current_value", data.CurrentValue,
		"target_value", data.TargetValue,
	)

	data.UpdateProgress()
	data.UpdateStatus()

	logger.Debug("Creating key result record",
		"key_result_id", data.ID,
		"title", data.Title,
		"progress", data.Progress,
		"status", data.Status,
	)

	created, err := k.repo.Create(&data)
	if err != nil {
		logger.Error("Failed to create key result record",
			"key_result_id", data.ID,
			"title", data.Title,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create Key Result: %w", err)
	}

	logger.Info("Key result created successfully",
		"key_result_id", created.ID,
		"title", created.Title,
		"objective_id", created.ObjectiveID,
		"progress", created.Progress,
		"status", created.Status,
	)

	return created, nil
}

func (k *keyResultService) GetData(identifier, id string) (*models.KeyResult, error) {
	logger.Debug("Retrieving key result",
		"identifier_type", identifier,
		"identifier_value", id,
	)

	res, err := k.repo.GetByIdentifier(identifier, id)
	if err != nil {
		logger.Error("Failed to retrieve key result",
			"identifier_type", identifier,
			"identifier_value", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get data: %v", err)
	}

	logger.Info("Key result retrieved successfully",
		"key_result_id", res.ID,
		"title", res.Title,
		"identifier_type", identifier,
		"identifier_value", id,
	)

	return res, nil
}

func (k *keyResultService) UpdateKeyResult(req dto.UpdateKeyResultRequest) (*models.KeyResult, error) {
	logger.Info("Key result update started",
		"key_result_id", req.ID,
		"title", req.Title,
		"current_value", req.CurrentValue,
		"target_value", req.TargetValue,
	)

	if err := k.validator.Struct(req); err != nil {
		logger.Error("Key result update failed - validation error",
			"key_result_id", req.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	logger.Debug("Retrieving existing key result for update",
		"key_result_id", req.ID,
	)

	existing, err := k.repo.GetByIdentifier("id", req.ID)
	if err != nil {
		logger.Error("Failed to find key result for update",
			"key_result_id", req.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to find key result: %w", err)
	}

	oldProgress := existing.Progress
	oldStatus := existing.Status

	logger.Debug("Applying updates to key result",
		"key_result_id", req.ID,
		"old_title", existing.Title,
		"new_title", req.Title,
		"old_current_value", existing.CurrentValue,
		"new_current_value", req.CurrentValue,
	)

	existing.Title = req.Title
	existing.Description = req.Description
	existing.CurrentValue = req.CurrentValue
	existing.TargetValue = req.TargetValue
	existing.MetricType = req.MetricType
	existing.AssigneeType = req.AssigneeType
	existing.AssigneeID = req.AssigneeID
	existing.StartDate = req.StartDate
	existing.DueDate = req.DueDate

	logger.Debug("Validating updated metric values",
		"key_result_id", existing.ID,
		"metric_type", existing.MetricType,
		"current_value", existing.CurrentValue,
		"target_value", existing.TargetValue,
	)

	if err := validation.ValidateMetricValues(existing); err != nil {
		logger.Error("Key result update failed - metric validation error",
			"key_result_id", existing.ID,
			"metric_type", existing.MetricType,
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug("Recalculating progress and status",
		"key_result_id", existing.ID,
		"old_progress", oldProgress,
		"old_status", oldStatus,
	)

	existing.UpdateProgress()
	existing.UpdateStatus()

	logger.Debug("Progress and status recalculated",
		"key_result_id", existing.ID,
		"new_progress", existing.Progress,
		"new_status", existing.Status,
		"progress_changed", oldProgress != existing.Progress,
		"status_changed", oldStatus != existing.Status,
	)

	updatedData, err := k.repo.Update(existing)
	if err != nil {
		logger.Error("Failed to update key result",
			"key_result_id", existing.ID,
			"title", existing.Title,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to update key Result: %v", err)
	}

	logger.Info("Key result updated successfully",
		"key_result_id", updatedData.ID,
		"title", updatedData.Title,
		"progress", updatedData.Progress,
		"status", updatedData.Status,
	)

	return updatedData, nil
}

func (k *keyResultService) DeleteKeyResult(id string) error {
	logger.Info("Key result deletion started",
		"key_result_id", id,
	)

	if err := k.repo.Delete(id); err != nil {
		logger.Error("Failed to delete key result",
			"key_result_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to delete key Result: %v", err)
	}

	logger.Info("Key result deleted successfully",
		"key_result_id", id,
	)

	return nil
}

func (k *keyResultService) ListData(identifier, objId string) ([]models.KeyResult, error) {
	logger.Debug("Retrieving key results by identifier",
		"identifier_type", identifier,
		"identifier_value", objId,
	)

	keys, err := k.repo.ListByIdentifier(identifier, objId)
	if err != nil {
		logger.Error("Failed to retrieve key results by identifier",
			"identifier_type", identifier,
			"identifier_value", objId,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to list data: %v", err)
	}

	logger.Info("Key results retrieved successfully",
		"identifier_type", identifier,
		"identifier_value", objId,
		"key_result_count", len(keys),
	)

	return keys, nil
}

// func (k *keyResultService) ListAssigneeKeyResults(identifier, userId string) (*models.KeyResult, error) {
// 	assignee, err := k.repo.GetByIdentifier(identifier, userId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get Assignee's Key Results: %v", err)
// 	}

// 	return assignee, nil
// }
