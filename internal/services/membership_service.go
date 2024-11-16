package services

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MembershipService interface {
	CreateMembership(r dto.CreateMembershipRequest) (*models.Membership, error)
	GetMembership(ident, id string) (*models.Membership, error)
	DeleteMembership(id string) error
	UpdateMembership(r dto.UpdateMembershipRequest) (*models.Membership, error)
	GetCompanyMembers(companyID string) ([]models.Membership, error)
	UpdateMembershipRole(id string, role models.RoleType) error
	UpdateMembershipStatus(id string, status models.StatusType) error
}

type membershipService struct {
	repo      repositories.MembershipRepository
	validator *validator.Validate
}

func NewMembershipService(repo repositories.MembershipRepository, validator *validator.Validate) MembershipService {
	return &membershipService{
		repo:      repo,
		validator: validator,
	}
}

func (m *membershipService) CreateMembership(r dto.CreateMembershipRequest) (*models.Membership, error) {
	if err := m.validator.Struct(r); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	membership := models.Membership{
		ID:        uuid.NewString(),
		UserID:    r.UserID,
		CompanyID: r.CompanyID,
		Role:      r.Role,
		Status:    models.StatusActive,
	}

	created, err := m.repo.Create(&membership)
	if err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return created, nil
}

func (m *membershipService) GetMembership(ident, id string) (*models.Membership, error) {
	membership, err := m.repo.GetByIdentifier(ident, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}
	return membership, nil
}

func (m *membershipService) UpdateMembership(r dto.UpdateMembershipRequest) (*models.Membership, error) {
	if err := m.validator.Struct(r); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// First get existing membership
	existing, err := m.repo.GetByIdentifier("id", r.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}

	// Update fields
	existing.Role = r.Role
	existing.Status = r.Status

	updated, err := m.repo.Update(existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update membership: %w", err)
	}

	return updated, nil
}

func (m *membershipService) DeleteMembership(id string) error {
	// Check if it's the last admin
	if err := m.validateDeletion(id); err != nil {
		return err
	}

	if err := m.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete membership: %w", err)
	}

	return nil
}

func (m *membershipService) GetCompanyMembers(companyID string) ([]models.Membership, error) {
	var memberships []models.Membership
	result := m.repo.GetDB().Where("company_id = ?", companyID).Find(&memberships)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get company members: %w", result.Error)
	}
	return memberships, nil
}

func (m *membershipService) UpdateMembershipRole(id string, role models.RoleType) error {
	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		return fmt.Errorf("failed to find membership: %w", err)
	}

	membership.Role = role
	_, err = m.repo.Update(membership)
	if err != nil {
		return fmt.Errorf("failed to update membership role: %w", err)
	}

	return nil
}

func (m *membershipService) UpdateMembershipStatus(id string, status models.StatusType) error {
	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		return fmt.Errorf("failed to find membership: %w", err)
	}

	membership.Status = status
	_, err = m.repo.Update(membership)
	if err != nil {
		return fmt.Errorf("failed to update membership status: %w", err)
	}

	return nil
}

// validateDeletion checks if the membership can be safely deleted
func (m *membershipService) validateDeletion(id string) error {
	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		return fmt.Errorf("failed to find membership: %w", err)
	}

	// If the membership is not an admin, we can safely delete it
	if membership.Role != models.RoleAdmin {
		return nil
	}

	// Count remaining active admins
	var adminCount int64
	result := m.repo.GetDB().Model(&models.Membership{}).
		Where("company_id = ? AND role = ? AND status = ? AND id != ?",
			membership.CompanyID, models.RoleAdmin, models.StatusActive, id).
		Count(&adminCount)

	if result.Error != nil {
		return fmt.Errorf("failed to count admins: %w", result.Error)
	}

	if adminCount == 0 {
		return fmt.Errorf("cannot delete the last admin of the company")
	}

	return nil
}
