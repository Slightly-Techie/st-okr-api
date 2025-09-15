package models

import "time"

type ObjectiveStatus string
type ObjectiveType string

const (
	ObjectiveStatusDraft      ObjectiveStatus = "draft"
	ObjectiveStatusActive     ObjectiveStatus = "active"
	ObjectiveStatusCompleted  ObjectiveStatus = "completed"
	ObjectiveStatusArchived   ObjectiveStatus = "archived"
	ObjectiveStatusOnHold     ObjectiveStatus = "on_hold"
)

const (
	ObjectiveTypeCompany ObjectiveType = "company"
	ObjectiveTypeTeam    ObjectiveType = "team"
)

type Objective struct {
	ID           string          `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Title        string          `gorm:"column:title;not null" json:"title,omitempty"`
	Description  string          `gorm:"column:description" json:"description,omitempty"`
	Type         ObjectiveType   `gorm:"column:type;type:varchar(50);not null;default:'team'" json:"type,omitempty" validate:"oneof=company team"`
	OwnerID      string          `gorm:"column:owner_id;not null;index" json:"owner_id,omitempty"`
	CompanyID    string          `gorm:"column:company_id;not null;index" json:"company_id,omitempty"`
	TeamID       *string         `gorm:"column:team_id;index" json:"team_id,omitempty"`
	Status       ObjectiveStatus `gorm:"column:status;type:varchar(50);default:'draft'" json:"status,omitempty" validate:"oneof=draft active completed archived on_hold"`
	StartDate    time.Time       `gorm:"column:start_date;not null" json:"start_date,omitempty"`
	EndDate      time.Time       `gorm:"column:end_date;not null" json:"end_date,omitempty"`
	Progress     float64         `gorm:"column:progress;default:0" json:"progress,omitempty"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime" json:"updated_at,omitempty"`
	
	// Relations
	KeyResults   []KeyResult     `gorm:"foreignKey:ObjectiveID" json:"key_results,omitempty"`
	Owner        User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Company      Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Team         *Team           `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (o *Objective) UpdateProgress() {
	if len(o.KeyResults) == 0 {
		o.Progress = 0
		return
	}

	totalProgress := 0.0
	for _, kr := range o.KeyResults {
		totalProgress += kr.Progress
	}
	o.Progress = totalProgress / float64(len(o.KeyResults))
}

func (o *Objective) UpdateStatus() {
	now := time.Now()
	
	switch {
	case o.Progress == 100:
		o.Status = ObjectiveStatusCompleted
	case now.After(o.EndDate) && o.Progress < 100:
		o.Status = ObjectiveStatusArchived
	case now.Before(o.StartDate):
		o.Status = ObjectiveStatusDraft
	case o.Progress > 0:
		o.Status = ObjectiveStatusActive
	}
}