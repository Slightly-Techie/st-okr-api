package models

import (
	"time"
)

// RoleType defines the available membership roles
type RoleType string

const (
	RoleAdmin  RoleType = "admin"
	RoleMember RoleType = "member"
	RoleViewer RoleType = "viewer"
)

// StatusType defines the available membership statuses
type StatusType string

const (
	StatusActive    StatusType = "active"
	StatusInactive  StatusType = "inactive"
	StatusSuspended StatusType = "suspended"
)

// Membership represents a user's membership in an organization or group
type Membership struct {
	ID        string     `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	UserID    string     `gorm:"column:user_id;not null;index" json:"user_id,omitempty"`
	CompanyID string     `gorm:"column:company_id;not null;index" json:"company_id,omitempty"`
	Role      RoleType   `gorm:"column:role;type:varchar(50);not null;default:'member'" json:"role,omitempty" validate:"required,oneof=admin member viewer"`
	Status    StatusType `gorm:"column:status;type:varchar(50);not null;default:'active'" json:"status,omitempty" validate:"required,oneof=active inactive suspended"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime" json:"updated_at,omitempty"`
}
