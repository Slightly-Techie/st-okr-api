package models

import "time"

type User struct {
	ID         string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	ProviderID string    `gorm:"column:provider_id;not null" json:"provider_id,omitempty"`
	FirstName  string    `gorm:"column:first_name;not null" json:"first_name,omitempty"`
	LastName   string    `gorm:"column:last_name;not null" json:"last_name,omitempty"`
	UserName   string    `gorm:"column:user_name;not null" json:"user_name,omitempty"`
	Email      string    `gorm:"column:email;unique;not null" json:"email,omitempty"`
	AvatarURL  string    `gorm:"column:avatar_url;unique;not null" json:"avatar_url,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}