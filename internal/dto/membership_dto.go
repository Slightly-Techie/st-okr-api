package dto

import "github.com/Slightly-Techie/st-okr-api/internal/models"

type CreateMembershipRequest struct {
	UserID    string           `json:"user_id" validate:"required,uuid"`
	CompanyID string           `json:"company_id" validate:"required,uuid"`
	Role      models.RoleType  `json:"role" validate:"required,oneof=admin member viewer"`
}

type UpdateMembershipRequest struct {
	ID     string             `json:"id" validate:"required,uuid"`
	Role   models.RoleType    `json:"role" validate:"required,oneof=admin member viewer"`
	Status models.StatusType  `json:"status" validate:"required,oneof=active inactive suspended"`
}