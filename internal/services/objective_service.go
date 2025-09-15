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
	if err := s.validator.Struct(req); err != nil {
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

	created, err := s.repo.Create(&objective)
	if err != nil {
		return nil, fmt.Errorf("failed to create objective: %w", err)
	}

	return created, nil
}

func (s *objectiveService) GetObjective(identifier, id string) (*models.Objective, error) {
	objective, err := s.repo.GetByIdentifier(identifier, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get objective: %v", err)
	}

	return objective, nil
}

func (s *objectiveService) GetObjectiveWithKeyResults(id string) (*dto.ObjectiveResponse, error) {
	objective, err := s.repo.GetWithKeyResults(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get objective with key results: %v", err)
	}

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

	return response, nil
}

func (s *objectiveService) UpdateObjective(req dto.UpdateObjectiveRequest) (*models.Objective, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	existing, err := s.repo.GetByIdentifier("id", req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find objective: %w", err)
	}

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
		return nil, fmt.Errorf("failed to update objective: %v", err)
	}

	return updated, nil
}

func (s *objectiveService) DeleteObjective(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete objective: %v", err)
	}
	return nil
}

func (s *objectiveService) ListObjectivesByCompany(companyID string) ([]dto.ObjectiveListResponse, error) {
	objectives, err := s.repo.ListByCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list objectives by company: %v", err)
	}

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) ListObjectivesByTeam(teamID string) ([]dto.ObjectiveListResponse, error) {
	objectives, err := s.repo.ListByTeam(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list objectives by team: %v", err)
	}

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) ListObjectivesByOwner(ownerID string) ([]dto.ObjectiveListResponse, error) {
	objectives, err := s.repo.ListByOwner(ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list objectives by owner: %v", err)
	}

	return s.mapToListResponse(objectives), nil
}

func (s *objectiveService) UpdateObjectiveProgress(objectiveID string) error {
	objective, err := s.repo.GetWithKeyResults(objectiveID)
	if err != nil {
		return fmt.Errorf("failed to get objective: %v", err)
	}

	objective.UpdateProgress()
	objective.UpdateStatus()

	_, err = s.repo.Update(objective)
	if err != nil {
		return fmt.Errorf("failed to update objective progress: %v", err)
	}

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