package models

import "time"

type Objectives struct {
	ID        string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Title     string    `gorm:"column:title;not null" json:"title,omitempty"`
	CreatorID string    `gorm:"column:creator_id;not null" json:"creator,omitempty"`
	Deadline  string    `gorm:"column:deadline;not null" json:"deadline,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}
