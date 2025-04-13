package models

import (
	"time"
)

type MetricType string
type KeyResultProgressStatus string
type AssigneeType string

const (
	MetricTypeNumeric    MetricType = "numeric"
	MetricTypePercentage MetricType = "percentage"
	MetricTypeBinary     MetricType = "binary"
	MetrictTypeCurrency  MetricType = "currency"
)

const (
	StatusNotStarted KeyResultProgressStatus = "not_started"
	StatusInProgress KeyResultProgressStatus = "on_track"
	StatusRisk       KeyResultProgressStatus = "at_risk"
	StatusBehind     KeyResultProgressStatus = "behind"
	StatusCompleted  KeyResultProgressStatus = "completed"
)

const (
	AssigneeTypeIndividual AssigneeType = "individual"
	AssigneeTypeTeam       AssigneeType = "team"
)

type KeyResult struct {
	ID           string                  `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	Title        string                  `gorm:"column:title;not null" json:"name,omitempty"`
	Description  string                  `gorm:"column:description;not null" json:"description,omitempty"`
	ObjectiveID  string                  `gorm:"column:objective_id;not null;index" json:"objective_id,omitempty"`
	MetricType   MetricType              `gorm:"column:metric_type;type:varchar(50);not null;default:'percentage'" json:"metric_type,omitempty" validate:"required,metric_type"`
	TargetValue  float64                 `gorm:"column:target_value;not null" json:"target_value"`
	CurrentValue float64                 `gorm:"column:current_value;" json:"current_value"`
	Progress     float64                 `gorm:"column:progress;default:0" json:"progress_percentage"`
	Status       KeyResultProgressStatus `gorm:"column:status;type:varchar(50);default:'not_started'" json:"status,omitempty" validate:"oneof=not_started on_track at_risk behind completed"`
	AssigneeType AssigneeType            `gorm:"column:assignee_type;type:varchar(50);not null;default:'team'" json:"assignee_type,omitempty" validate:"assignee_type"`
	AssigneeID   string                  `gorm:"column:assignee_id;not null;index" json:"assignee_id,omitempty"`
	StartDate    time.Time               `gorm:"column:start_date; not null" json:"start_date,omitempty"`
	DueDate      time.Time               `gorm:"column:due_date; not null" json:"due_date,omitempty" validate:"due_date"`
	CreatedAt    time.Time               `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt    time.Time               `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime" json:"updated_at,omitempty"`
	// UpdatedBy []string  `gorm:"column:updated_by;type:jsonb;index" json:"user_id,omitempty"`
}

// func (k *KeyResult) UpdateProgress() {
// 	if k.TargetValue == 0 {
// 		k.Progress = 0
// 	}
// 	k.Progress = (k.CurrentValue / k.TargetValue) * 100
// }

func (k *KeyResult) UpdateProgress() {
	switch k.MetricType {
	case MetricTypeBinary:
		if k.CurrentValue == 1 {
			k.TargetValue = 1
			k.Progress = 100
		}
	case MetricTypeNumeric, MetrictTypeCurrency, MetricTypePercentage:
		if k.TargetValue == 0 {
			k.Progress = 0
		}
		if k.TargetValue > 0 {
			k.Progress = (k.CurrentValue / k.TargetValue) * 100
		}

	}
}

func (k *KeyResult) UpdateStatus() {
	now := time.Now()

	switch {
	case k.Progress == 0 && now.Before(k.StartDate):
		k.Status = "not_started"

	case k.Progress == 100:
		k.Status = "completed"

	case k.Progress > 0 && k.Progress < 100:
		if now.After(k.DueDate) || (k.DueDate.Sub(now).Hours() < 24*7 && k.Progress < 50) {
			k.Status = "at_risk"
		} else {
			k.Status = "on_track"
		}
	default:
		k.Status = "behind"
	}
}
