package models

import "time"

type Company struct {
	ID        string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Name      string    `gorm:"column:name;not null" json:"name,omitempty"`
	Code      string    `gorm:"column:company_code;not null" json:"company_code,omitempty"`
	CreatorID string    `gorm:"column:creator_id;not null" json:"creator_id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}
