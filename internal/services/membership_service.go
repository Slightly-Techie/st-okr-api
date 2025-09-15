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
	logger.Info("Membership creation started",
		"user_id", r.UserID,
		"company_id", r.CompanyID,
		"role", r.Role,
	)

	if err := m.validator.Struct(r); err != nil {
		logger.Error("Membership creation failed - validation error",
			"user_id", r.UserID,
			"company_id", r.CompanyID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	membership := models.Membership{
		ID:        uuid.NewString(),
		UserID:    r.UserID,
		CompanyID: r.CompanyID,
		Role:      r.Role,
		Status:    models.StatusActive,
	}

	logger.Debug("Creating membership record",
		"membership_id", membership.ID,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
		"role", membership.Role,
		"status", membership.Status,
	)

	created, err := m.repo.Create(&membership)
	if err != nil {
		logger.Error("Failed to create membership record",
			"membership_id", membership.ID,
			"user_id", membership.UserID,
			"company_id", membership.CompanyID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	logger.Info("Membership created successfully",
		"membership_id", created.ID,
		"user_id", created.UserID,
		"company_id", created.CompanyID,
		"role", created.Role,
		"status", created.Status,
	)

	return created, nil
}

func (m *membershipService) GetMembership(ident, id string) (*models.Membership, error) {
	logger.Debug("Retrieving membership",
		"identifier_type", ident,
		"identifier_value", id,
	)

	membership, err := m.repo.GetByIdentifier(ident, id)
	if err != nil {
		logger.Error("Failed to retrieve membership",
			"identifier_type", ident,
			"identifier_value", id,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	logger.Info("Membership retrieved successfully",
		"membership_id", membership.ID,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
		"role", membership.Role,
		"identifier_type", ident,
		"identifier_value", id,
	)

	return membership, nil
}

func (m *membershipService) UpdateMembership(r dto.UpdateMembershipRequest) (*models.Membership, error) {
	logger.Info("Membership update started",
		"membership_id", r.ID,
		"new_role", r.Role,
		"new_status", r.Status,
	)

	if err := m.validator.Struct(r); err != nil {
		logger.Error("Membership update failed - validation error",
			"membership_id", r.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	logger.Debug("Retrieving existing membership for update",
		"membership_id", r.ID,
	)

	// First get existing membership
	existing, err := m.repo.GetByIdentifier("id", r.ID)
	if err != nil {
		logger.Error("Failed to find membership for update",
			"membership_id", r.ID,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}

	logger.Debug("Applying updates to membership",
		"membership_id", r.ID,
		"old_role", existing.Role,
		"new_role", r.Role,
		"old_status", existing.Status,
		"new_status", r.Status,
	)

	// Update fields
	existing.Role = r.Role
	existing.Status = r.Status

	updated, err := m.repo.Update(existing)
	if err != nil {
		logger.Error("Failed to update membership",
			"membership_id", r.ID,
			"role", existing.Role,
			"status", existing.Status,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to update membership: %w", err)
	}

	logger.Info("Membership updated successfully",
		"membership_id", updated.ID,
		"user_id", updated.UserID,
		"company_id", updated.CompanyID,
		"role", updated.Role,
		"status", updated.Status,
	)

	return updated, nil
}

func (m *membershipService) DeleteMembership(id string) error {
	logger.Info("Membership deletion started",
		"membership_id", id,
	)

	logger.Debug("Validating membership deletion",
		"membership_id", id,
	)

	// Check if it's the last admin
	if err := m.validateDeletion(id); err != nil {
		logger.Warn("Membership deletion validation failed",
			"membership_id", id,
			"error", err.Error(),
		)
		return err
	}

	if err := m.repo.Delete(id); err != nil {
		logger.Error("Failed to delete membership",
			"membership_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to delete membership: %w", err)
	}

	logger.Info("Membership deleted successfully",
		"membership_id", id,
	)

	return nil
}

func (m *membershipService) GetCompanyMembers(companyID string) ([]models.Membership, error) {
	logger.Debug("Retrieving company members",
		"company_id", companyID,
	)

	var memberships []models.Membership
	result := m.repo.GetDB().Where("company_id = ?", companyID).Find(&memberships)
	if result.Error != nil {
		logger.Error("Failed to retrieve company members",
			"company_id", companyID,
			"error", result.Error.Error(),
		)
		return nil, fmt.Errorf("failed to get company members: %w", result.Error)
	}

	logger.Info("Company members retrieved successfully",
		"company_id", companyID,
		"member_count", len(memberships),
	)

	return memberships, nil
}

func (m *membershipService) UpdateMembershipRole(id string, role models.RoleType) error {
	logger.Info("Membership role update started",
		"membership_id", id,
		"new_role", role,
	)

	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		logger.Error("Failed to find membership for role update",
			"membership_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to find membership: %w", err)
	}

	oldRole := membership.Role

	logger.Debug("Updating membership role",
		"membership_id", id,
		"old_role", oldRole,
		"new_role", role,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
	)

	membership.Role = role
	_, err = m.repo.Update(membership)
	if err != nil {
		logger.Error("Failed to update membership role",
			"membership_id", id,
			"old_role", oldRole,
			"new_role", role,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to update membership role: %w", err)
	}

	logger.Info("Membership role updated successfully",
		"membership_id", id,
		"old_role", oldRole,
		"new_role", role,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
	)

	return nil
}

func (m *membershipService) UpdateMembershipStatus(id string, status models.StatusType) error {
	logger.Info("Membership status update started",
		"membership_id", id,
		"new_status", status,
	)

	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		logger.Error("Failed to find membership for status update",
			"membership_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to find membership: %w", err)
	}

	oldStatus := membership.Status

	logger.Debug("Updating membership status",
		"membership_id", id,
		"old_status", oldStatus,
		"new_status", status,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
	)

	membership.Status = status
	_, err = m.repo.Update(membership)
	if err != nil {
		logger.Error("Failed to update membership status",
			"membership_id", id,
			"old_status", oldStatus,
			"new_status", status,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to update membership status: %w", err)
	}

	logger.Info("Membership status updated successfully",
		"membership_id", id,
		"old_status", oldStatus,
		"new_status", status,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
	)

	return nil
}

// validateDeletion checks if the membership can be safely deleted
func (m *membershipService) validateDeletion(id string) error {
	logger.Debug("Validating membership deletion eligibility",
		"membership_id", id,
	)

	membership, err := m.repo.GetByIdentifier("id", id)
	if err != nil {
		logger.Error("Failed to find membership for deletion validation",
			"membership_id", id,
			"error", err.Error(),
		)
		return fmt.Errorf("failed to find membership: %w", err)
	}

	logger.Debug("Checking membership role for deletion validation",
		"membership_id", id,
		"role", membership.Role,
		"user_id", membership.UserID,
		"company_id", membership.CompanyID,
	)

	// If the membership is not an admin, we can safely delete it
	if membership.Role != models.RoleAdmin {
		logger.Debug("Non-admin membership can be safely deleted",
			"membership_id", id,
			"role", membership.Role,
		)
		return nil
	}

	logger.Debug("Admin membership found, checking for remaining admins",
		"membership_id", id,
		"company_id", membership.CompanyID,
	)

	// Count remaining active admins
	var adminCount int64
	result := m.repo.GetDB().Model(&models.Membership{}).
		Where("company_id = ? AND role = ? AND status = ? AND id != ?",
			membership.CompanyID, models.RoleAdmin, models.StatusActive, id).
		Count(&adminCount)

	if result.Error != nil {
		logger.Error("Failed to count remaining admins",
			"membership_id", id,
			"company_id", membership.CompanyID,
			"error", result.Error.Error(),
		)
		return fmt.Errorf("failed to count admins: %w", result.Error)
	}

	logger.Debug("Remaining admin count checked",
		"membership_id", id,
		"company_id", membership.CompanyID,
		"remaining_admin_count", adminCount,
	)

	if adminCount == 0 {
		logger.Warn("Cannot delete last admin of company",
			"membership_id", id,
			"company_id", membership.CompanyID,
			"user_id", membership.UserID,
		)
		return fmt.Errorf("cannot delete the last admin of the company")
	}

	logger.Debug("Admin deletion validated successfully",
		"membership_id", id,
		"company_id", membership.CompanyID,
		"remaining_admin_count", adminCount,
	)

	return nil
}
