package models

import "time"

type Team struct {
	ID          string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Name        string    `gorm:"column:name;not null" json:"name,omitempty"`
	CompanyID   string    `gorm:"column:company_id;not null;index" json:"company_id,omitempty"`
	Description string    `gorm:"column:description;not null" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

type TeamMember struct {
	ID     string `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	UserID string `gorm:"column:user_id;not null;index" json:"user_id,omitempty"`
	TeamID string `gorm:"column:team_id;not null;index" json:"team_id,omitempty"`
}
