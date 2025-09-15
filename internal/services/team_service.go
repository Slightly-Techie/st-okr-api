package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
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
	logger.Info("Team creation started",
		"team_name", t.Name,
		"company_id", t.CompanyID,
		"description", t.Description,
	)

	if err := r.validator.Struct(t); err != nil {
		logger.Error("Team creation failed - validation error",
			"team_name", t.Name,
			"company_id", t.CompanyID,
			"error", err.Error(),
		)
		return nil, err
	}

	team := models.Team{
		ID:          uuid.NewString(),
		Name:        t.Name,
		CompanyID:   t.CompanyID,
		Description: t.Description,
	}

	logger.Debug("Creating team record",
		"team_id", team.ID,
		"team_name", team.Name,
		"company_id", team.CompanyID,
	)

	created, err := r.repo.CreateTeam(&team)
	if err != nil {
		logger.Error("Failed to create team record",
			"team_id", team.ID,
			"team_name", team.Name,
			"company_id", team.CompanyID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	logger.Info("Team created successfully",
		"team_id", created.ID,
		"team_name", created.Name,
		"company_id", created.CompanyID,
	)

	return created, nil
}

func (r *teamService) GetTeam(identifier, id string) (*models.Team, error) {
	logger.Debug("Retrieving team",
		"identifier_type", identifier,
		"identifier_value", id,
	)

	team, err := r.repo.GetByIdentifier(identifier, id)
	if err != nil {
		logger.Error("Failed to retrieve team",
			"identifier_type", identifier,
			"identifier_value", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get team: %v", err)
	}

	logger.Info("Team retrieved successfully",
		"team_id", team.ID,
		"team_name", team.Name,
		"identifier_type", identifier,
		"identifier_value", id,
	)

	return team, nil
}

func (r *teamService) UpdateTeam(t dto.UpdateTeamRequest) (*models.Team, error) {
	logger.Info("Team update started",
		"team_id", t.ID,
		"team_name", t.Name,
		"description", t.Description,
	)

	if err := r.validator.Struct(t); err != nil {
		logger.Error("Team update failed - validation error",
			"team_id", t.ID,
			"team_name", t.Name,
			"error", err.Error(),
		)
		return nil, err
	}

	team := models.Team{
		Name:        t.Name,
		Description: t.Description,
	}

	logger.Debug("Updating team record",
		"team_id", t.ID,
		"team_name", team.Name,
		"description", team.Description,
	)

	updatedTeam, err := r.repo.UpdateTeam(&team)
	if err != nil {
		logger.Error("Failed to update team",
			"team_id", t.ID,
			"team_name", team.Name,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to update team: %v", err)
	}

	logger.Info("Team updated successfully",
		"team_id", updatedTeam.ID,
		"team_name", updatedTeam.Name,
		"description", updatedTeam.Description,
	)

	return updatedTeam, nil
}

func (r *teamService) DeleteTeam(id string) error {
	logger.Info("Team deletion started",
		"team_id", id,
	)

	if err := r.repo.DeleteTeam(id); err != nil {
		logger.Error("Failed to delete team",
			"team_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to delete team: %v", err)
	}

	logger.Info("Team deleted successfully",
		"team_id", id,
	)

	return nil
}

func (r *teamService) AddMember(t *dto.TeamMemberRequest) (*models.TeamMember, error) {
	logger.Info("Team member addition started",
		"team_id", t.TeamID,
		"user_id", t.UserID,
	)

	if err := r.validator.Struct(t); err != nil {
		logger.Error("Team member addition failed - validation error",
			"team_id", t.TeamID,
			"user_id", t.UserID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("validation error: %v", err)
	}

	logger.Debug("Checking existing team membership",
		"team_id", t.TeamID,
		"user_id", t.UserID,
	)

	isMember, err := r.repo.IsMember(t.TeamID, t.UserID)
	if err != nil {
		logger.Error("Failed to check team membership",
			"team_id", t.TeamID,
			"user_id", t.UserID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to check user's team membership: %v", err)
	}

	if isMember {
		logger.Warn("User is already a team member",
			"team_id", t.TeamID,
			"user_id", t.UserID,
		)
		return nil, fmt.Errorf("user is already a member of the team")
	}

	teamMember := models.TeamMember{
		ID:     uuid.NewString(),
		UserID: t.UserID,
		TeamID: t.TeamID,
	}

	logger.Debug("Creating team membership record",
		"member_id", teamMember.ID,
		"team_id", teamMember.TeamID,
		"user_id", teamMember.UserID,
	)

	created, err := r.repo.AddTeamMember(&teamMember)
	if err != nil {
		logger.Error("Failed to create team membership",
			"member_id", teamMember.ID,
			"team_id", teamMember.TeamID,
			"user_id", teamMember.UserID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create team membership: %w", err)
	}

	logger.Info("Team member added successfully",
		"member_id", created.ID,
		"team_id", created.TeamID,
		"user_id", created.UserID,
	)

	return created, nil
}

func (r *teamService) ListMembers(identifier, teamID string) ([]models.TeamMember, error) {
	logger.Debug("Retrieving team members",
		"identifier_type", identifier,
		"team_id", teamID,
	)

	teamMembers, err := r.repo.GetTeamMembers(identifier, teamID)
	if err != nil {
		logger.Error("Failed to retrieve team members",
			"identifier_type", identifier,
			"team_id", teamID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}

	logger.Info("Team members retrieved successfully",
		"identifier_type", identifier,
		"team_id", teamID,
		"member_count", len(teamMembers),
	)

	return teamMembers, nil
}

func (r *teamService) RemoveMember(id string) error {
	logger.Info("Team member removal started",
		"member_id", id,
	)

	if err := r.repo.RemoveTeamMember(id); err != nil {
		logger.Error("Failed to remove team member",
			"member_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to remove team member: %v", err)
	}

	logger.Info("Team member removed successfully",
		"member_id", id,
	)

	return nil
}

// func (r *teamService) isTeamMember(t dto.TeamMemberRequest) (bool, error) {

// }
