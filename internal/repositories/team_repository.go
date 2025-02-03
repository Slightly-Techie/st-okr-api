package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrTeamNotFound       = errors.New("no Team exists with the provided credentials")
	ErrTeamDBOperation    = errors.New("database operation failed")
	ErrTeamMemberNotFound = errors.New("no Team member exists with the provided credentials")
)

type TeamRepository interface {
	GetDB() *gorm.DB
	GetByIdentifier(identifier, id string) (*models.Team, error)
	CreateTeam(team *models.Team) (*models.Team, error)
	UpdateTeam(team *models.Team) (*models.Team, error)
	DeleteTeam(id string) error

	AddTeamMember(member *models.TeamMember) (*models.TeamMember, error)
	RemoveTeamMember(id string) error
	GetTeamMembers(id string) ([]models.TeamMember, error)
	IsMember(teamID, userID string) (bool, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		db: db,
	}
}

func (r *teamRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *teamRepository) GetByIdentifier(identifier, id string) (*models.Team, error) {
	var team models.Team

	res := r.db.Where(identifier, id).First(&team)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		log.Printf("error getting team by identifier: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	return &team, nil
}

func (r *teamRepository) CreateTeam(team *models.Team) (*models.Team, error) {
	res := r.db.Create(team)
	if res.Error != nil {
		log.Printf("error creating team: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	return team, nil
}

func (r *teamRepository) UpdateTeam(team *models.Team) (*models.Team, error) {
	res := r.db.Save(team)
	if res.Error != nil {
		log.Printf("error updating team: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	return team, nil
}

func (r *teamRepository) DeleteTeam(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Team{})
	if res.Error != nil {
		log.Printf("error deleting team: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no team found with id: %s", id)
		return ErrTeamNotFound
	}
	return nil
}

func (r *teamRepository) AddTeamMember(member *models.TeamMember) (*models.TeamMember, error) {
	res := r.db.Create(member)
	if res.Error != nil {
		log.Printf("error adding team member: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	return member, nil
}

func (r *teamRepository) GetTeamMembers(id string) ([]models.TeamMember, error) {
	var members []models.TeamMember

	res := r.db.Where("id = ?", id).Find(&members)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTeamMemberNotFound
		}
		log.Printf("error getting team members: %v", res.Error)
		return nil, fmt.Errorf("%w: %v", ErrTeamMemberNotFound, res.Error)
	}
	return members, nil
}

func (r *teamRepository) RemoveTeamMember(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.TeamMember{})
	if res.Error != nil {
		log.Printf("error removing team member: %v", res.Error)
		return fmt.Errorf("%w: %v", ErrTeamDBOperation, res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("no team member found with id: %s", id)
		return ErrTeamMemberNotFound
	}
	return nil
}

func (r *teamRepository) IsMember(teamID, userID string) (bool, error) {
	var member models.TeamMember

	res := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Printf("error checking team membership: %v", res.Error)
		return false, fmt.Errorf("failed to check team membership: %v", res.Error)
	}
	return true, nil
}
