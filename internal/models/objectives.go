package models

import "time"

type Objective struct {
	ID        string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Title     string    `gorm:"column:title;not null" json:"title,omitempty"`
	Creator   string    `gorm:"column:creator;not null" json:"creator,omitempty"`
	Deadline  string    `gorm:"column:deadline;not null" json:"deadline,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}
