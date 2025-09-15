package services

import (
	"fmt"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ObjectiveService interface {
	CreateObjective(req dto.CreateObjectiveRequest) (*models.Objective, error)
	GetObjective(identifier, id string) (*models.Objective, error)
	GetObjectiveWithKeyResults(id string) (*dto.ObjectiveResponse, error)
	UpdateObjective(req dto.UpdateObjectiveRequest) (*models.Objective, error)
	DeleteObjective(id string) error
	ListObjectivesByCompany(companyID string) ([]dto.ObjectiveListResponse, error)
	ListObjectivesByTeam(teamID string) ([]dto.ObjectiveListResponse, error)
	ListObjectivesByOwner(ownerID string) ([]dto.ObjectiveListResponse, error)
	UpdateObjectiveProgress(objectiveID string) error
}

type objectiveService struct {
	repo      repositories.ObjectiveRepository
	validator *validator.Validate
}

func NewObjectiveService(repo repositories.ObjectiveRepository, validator *validator.Validate) ObjectiveService {
	return &objectiveService{
		repo:      repo,
		validator: validator,
	}
}

func (s *objectiveService) CreateObjective(req dto.CreateObjectiveRequest) (*models.Objective, error) {
	logger.Info("Objective creation started",
		"title", req.Title,
		"type", req.Type,
		"owner_id", req.OwnerID,
		"company_id", req.CompanyID,
		"team_id", req.TeamID,
	)

	if err := s.validator.Struct(req); err != nil {
		logger.Error("Objective creation failed - validation error",
			"title", req.Title,
			"owner_id", req.OwnerID,
			"error", err.Error(),
		)
		return nil, err
	}

	objective := models.Objective{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		OwnerID:     req.OwnerID,
		CompanyID:   req.CompanyID,
		TeamID:      req.TeamID,
		Status:      models.ObjectiveStatusDraft,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Progress:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	logger.Debug("Creating objective record",
		"objective_id", objective.ID,
		"title", objective.Title,
		"type", objective.Type,
		"status", objective.Status,
	)

	created, err := s.repo.Create(&objective)
	if err != nil {
		logger.Error("Failed to create objective record",
			"objective_id", objective.ID,
			"title", objective.Title,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create objective: %w", err)
	}

	logger.Info("Objective created successfully",
		"objective_id", created.ID,
		"title", created.Title,
		"type", created.Type,
		"owner_id", created.OwnerID,
	)

	return created, nil
}

func (s *objectiveService) GetObjective(identifier, id string) (*models.Objective, error) {
	logger.Debug("Retrieving objective",
		"identifier_type", identifier,
		"identifier_value", id,
	)

	objective, err := s.repo.GetByIdentifier(identifier, id)
	if err != nil {
		logger.Error("Failed to retrieve objective",
			"identifier_type", identifier,
			"identifier_value", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get objective: %v", err)
	}

	logger.Info("Objective retrieved successfully",
		"objective_id", objective.ID,
		"title", objective.Title,
		"identifier_type", identifier,
		"identifier_value", id,
	)

	return objective, nil
}

func (s *objectiveService) GetObjectiveWithKeyResults(id string) (*dto.ObjectiveResponse, error) {
	logger.Debug("Retrieving objective with key results",
		"objective_id", id,
	)

	objective, err := s.repo.GetWithKeyResults(id)
	if err != nil {
		logger.Error("Failed to retrieve objective with key results",
			"objective_id", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get objective with key results: %v", err)
	}

	logger.Debug("Mapping key results to response",
		"objective_id", objective.ID,
		"key_result_count", len(objective.KeyResults),
	)

	keyResults := make([]dto.KeyResultResponse, len(objective.KeyResults))
	for i, kr := range objective.KeyResults {
		keyResults[i] = dto.KeyResultResponse{
			ID:           kr.ID,
			ObjectiveID:  kr.ObjectiveID,
			Title:        kr.Title,
			Description:  kr.Description,
			AssigneeType: kr.AssigneeType,
			AssigneeID:   kr.AssigneeID,
			MetricType:   kr.MetricType,
			CurrentValue: kr.CurrentValue,
			TargetValue:  kr.TargetValue,
			Progress:     kr.Progress,
			StartDate:    kr.StartDate,
			DueDate:      kr.DueDate,
			Status:       kr.Status,
		}
	}

	response := &dto.ObjectiveResponse{
		ID:          objective.ID,
		Title:       objective.Title,
		Description: objective.Description,
		Type:        objective.Type,
		OwnerID:     objective.OwnerID,
		CompanyID:   objective.CompanyID,
		TeamID:      objective.TeamID,
		Status:      objective.Status,
		StartDate:   objective.StartDate,
		EndDate:     objective.EndDate,
		Progress:    objective.Progress,
		CreatedAt:   objective.CreatedAt,
		UpdatedAt:   objective.UpdatedAt,
		KeyResults:  keyResults,
	}

	logger.Info("Objective with key results retrieved successfully",
		"objective_id", objective.ID,
		"title", objective.Title,
		"key_result_count", len(keyResults),
	)

	return response, nil
}

func (s *objectiveService) UpdateObjective(req dto.UpdateObjectiveRequest) (*models.Objective, error) {
	logger.Info("Objective update started",
		"objective_id", req.ID,
		"title", req.Title,
		"status", req.Status,
	)

	if err := s.validator.Struct(req); err != nil {
		logger.Error("Objective update failed - validation error",
			"objective_id", req.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	logger.Debug("Retrieving existing objective for update",
		"objective_id", req.ID,
	)

	existing, err := s.repo.GetByIdentifier("id", req.ID)
	if err != nil {
		logger.Error("Failed to find objective for update",
			"objective_id", req.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to find objective: %w", err)
	}

	logger.Debug("Applying updates to objective",
		"objective_id", req.ID,
		"old_title", existing.Title,
		"new_title", req.Title,
		"old_status", existing.Status,
		"new_status", req.Status,
	)

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if !req.StartDate.IsZero() {
		existing.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		existing.EndDate = req.EndDate
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.repo.Update(existing)
	if err != nil {
		logger.Error("Failed to update objective",
			"objective_id", req.ID,
			"title", existing.Title,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to update objective: %v", err)
	}

	logger.Info("Objective updated successfully",
		"objective_id", updated.ID,
		"title", updated.Title,
		"status", updated.Status,
	)

	return updated, nil
}

func (s *objectiveService) DeleteObjective(id string) error {
	logger.Info("Objective deletion started",
		"objective_id", id,
	)

	if err := s.repo.Delete(id); err != nil {
		logger.Error("Failed to delete objective",
			"objective_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to delete objective: %v", err)
	}

	logger.Info("Objective deleted successfully",
		"objective_id", id,
	)

	return nil
}

func (s *objectiveService) ListObjectivesByCompany(companyID string) ([]dto.ObjectiveListResponse, error) {
	logger.Debug("Retrieving objectives by company",
		"company_id", companyID,
	)

	objectives, err := s.repo.ListByCompany(companyID)
	if err != nil {
		logger.Error("Failed to retrieve objectives by company",
			"company_id", companyID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to list objectives by company: %v", err)
	}

	logger.Info("Company objectives retrieved successfully",
		"company_id", companyID,
		"objective_count", len(objectives),
	)

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) ListObjectivesByTeam(teamID string) ([]dto.ObjectiveListResponse, error) {
	logger.Debug("Retrieving objectives by team",
		"team_id", teamID,
	)

	objectives, err := s.repo.ListByTeam(teamID)
	if err != nil {
		logger.Error("Failed to retrieve objectives by team",
			"team_id", teamID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to list objectives by team: %v", err)
	}

	logger.Info("Team objectives retrieved successfully",
		"team_id", teamID,
		"objective_count", len(objectives),
	)

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) ListObjectivesByOwner(ownerID string) ([]dto.ObjectiveListResponse, error) {
	logger.Debug("Retrieving objectives by owner",
		"owner_id", ownerID,
	)

	objectives, err := s.repo.ListByOwner(ownerID)
	if err != nil {
		logger.Error("Failed to retrieve objectives by owner",
			"owner_id", ownerID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to list objectives by owner: %v", err)
	}

	logger.Info("Owner objectives retrieved successfully",
		"owner_id", ownerID,
		"objective_count", len(objectives),
	)

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) UpdateObjectiveProgress(objectiveID string) error {
	logger.Info("Objective progress update started",
		"objective_id", objectiveID,
	)

	logger.Debug("Retrieving objective with key results for progress calculation",
		"objective_id", objectiveID,
	)

	objective, err := s.repo.GetWithKeyResults(objectiveID)
	if err != nil {
		logger.Error("Failed to get objective for progress update",
			"objective_id", objectiveID,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to get objective: %v", err)
	}

	oldProgress := objective.Progress
	oldStatus := objective.Status

	logger.Debug("Calculating objective progress and status",
		"objective_id", objectiveID,
		"old_progress", oldProgress,
		"old_status", oldStatus,
		"key_result_count", len(objective.KeyResults),
	)

	objective.UpdateProgress()
	objective.UpdateStatus()

	logger.Debug("Progress and status calculated",
		"objective_id", objectiveID,
		"new_progress", objective.Progress,
		"new_status", objective.Status,
		"progress_changed", oldProgress != objective.Progress,
		"status_changed", oldStatus != objective.Status,
	)

	_, err = s.repo.Update(objective)
	if err != nil {
		logger.Error("Failed to update objective progress",
			"objective_id", objectiveID,
			"progress", objective.Progress,
			"status", objective.Status,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to update objective progress: %v", err)
	}

	logger.Info("Objective progress updated successfully",
		"objective_id", objectiveID,
		"progress", objective.Progress,
		"status", objective.Status,
	)

	return nil
}

func (s *objectiveService) mapToListResponse(objectives []models.Objective) []dto.ObjectiveListResponse {
	response := make([]dto.ObjectiveListResponse, len(objectives))
	for i, obj := range objectives {
		response[i] = dto.ObjectiveListResponse{
			ID:              obj.ID,
			Title:           obj.Title,
			Description:     obj.Description,
			Type:            obj.Type,
			Status:          obj.Status,
			StartDate:       obj.StartDate,
			EndDate:         obj.EndDate,
			Progress:        obj.Progress,
			KeyResultsCount: len(obj.KeyResults),
		}
	}
	return response
}