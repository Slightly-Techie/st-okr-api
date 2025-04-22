package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(t dto.CreateTeamRequest) (*models.Team, error)
	GetTeam(identifier, id string) (*models.Team, error)
	UpdateTeam(t dto.UpdateTeamRequest) (*models.Team, error)
	DeleteTeam(id string) error

	// AddMember(teamID, userID string) (*models.TeamMember, error)
	AddMember(t *dto.TeamMemberRequest) (*models.TeamMember, error)
	ListMembers(identifier, teamID string) ([]models.TeamMember, error)
	RemoveMember(id string) error
	// isTeamMember(t dto.TeamMemberRequest) (bool, error)
}

type teamService struct {
	repo      repositories.TeamRepository
	validator *validator.Validate
}

func NewTeamService(repo repositories.TeamRepository, validator *validator.Validate) TeamService {
	return &teamService{
		repo:      repo,
		validator: validator,
	}
}

func (r *teamService) CreateTeam(t dto.CreateTeamRequest) (*models.Team, error) {
	if err := r.validator.Struct(t); err != nil {
		return nil, err
	}

	team := models.Team{
		ID:          uuid.NewString(),
		Name:        t.Name,
		CompanyID:   t.CompanyID,
		Description: t.Description,
	}

	created, err := r.repo.CreateTeam(&team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return created, nil
}

func (r *teamService) GetTeam(identifier, id string) (*models.Team, error) {
	team, err := r.repo.GetByIdentifier(identifier, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %v", err)
	}
	return team, nil
}

func (r *teamService) UpdateTeam(t dto.UpdateTeamRequest) (*models.Team, error) {
	if err := r.validator.Struct(t); err != nil {
		return nil, err
	}

	team := models.Team{
		Name:        t.Name,
		Description: t.Description,
	}

	updatedTeam, err := r.repo.UpdateTeam(&team)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %v", err)
	}

	return updatedTeam, nil
}

func (r *teamService) DeleteTeam(id string) error {
	if err := r.repo.DeleteTeam(id); err != nil {
		return fmt.Errorf("failed to delete team: %v", err)
	}
	return nil
}

func (r *teamService) AddMember(t *dto.TeamMemberRequest) (*models.TeamMember, error) {
	if err := r.validator.Struct(t); err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}

	isMember, err := r.repo.IsMember(t.TeamID, t.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user's team membership: %v", err)
	}

	if isMember {
		return nil, fmt.Errorf("user is already a member of the team")
	}

	teamMember := models.TeamMember{
		ID:     uuid.NewString(),
		UserID: t.UserID,
		TeamID: t.TeamID,
	}

	created, err := r.repo.AddTeamMember(&teamMember)
	if err != nil {
		return nil, fmt.Errorf("failed to create team membership: %w", err)
	}

	return created, nil
}

func (r *teamService) ListMembers(identifier, teamID string) ([]models.TeamMember, error) {
	teamMembers, err := r.repo.GetTeamMembers(identifier, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}

	return teamMembers, nil
}

func (r *teamService) RemoveMember(id string) error {
	if err := r.repo.RemoveTeamMember(id); err != nil {
		return fmt.Errorf("failed to remove team member: %v", err)
	}

	return nil
}

// func (r *teamService) isTeamMember(t dto.TeamMemberRequest) (bool, error) {

// }
