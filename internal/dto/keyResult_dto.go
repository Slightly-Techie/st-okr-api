package dto

import (
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
)

type CreateKeyResultRequest struct {
	ObjectiveID  string              `json:"objective_id" validate:"required,uuid"`
	Title        string              `json:"title" validate:"required"`
	Description  string              `json:"description"`
	MetricType   models.MetricType   `json:"metric_type" validate:"required,oneof=numeric percentage binary currency"`
	CurrentValue float64             `json:"current_value" validate:"required"`
	TargetValue  float64             `json:"target_value" validate:"required"`
	AssigneeType models.AssigneeType `json:"assignee_type" validate:"required,oneof=individual team"`
	AssigneeID   string              `json:"assignee_id" validate:"required,uuid"`
	StartDate    time.Time           `json:"start_date" validate:"required"`
	DueDate      time.Time           `json:"due_date" validate:"required"`
}

type UpdateKeyResultRequest struct {
	Title        string                         `json:"title"`
	Description  string                         `json:"description"`
	CurrentValue float64                        `json:"current_value"`
	Progress     float64                        `json:"progress" validate:"min=0,max=100"`
	Status       models.KeyResultProgressStatus `json:"status" validate:"oneof=not_started on_track at_risk behind completed"`
	AssigneeType models.AssigneeType            `json:"assignee_type" validate:"oneof=individual team"`
	AssigneeID   string                         `json:"assignee_id" validate:"uuid"`
	DueDate      time.Time                      `json:"due_date" validate:"future"`
}

type KeyResultResponse struct {
	ID           string                         `json:"id"`
	ObjectiveID  string                         `json:"objective_id"`
	Title        string                         `json:"title"`
	Description  string                         `json:"description"`
	AssigneeType models.AssigneeType            `json:"assignee_type"`
	AssigneeID   string                         `json:"assignee_id"`
	MetricType   models.MetricType              `json:"metric_type"`
	CurrentValue float64                        `json:"current_value"`
	TargetValue  float64                        `json:"target_value"`
	Progress     float64                        `json:"progress"`
	StartDate    time.Time                      `json:"start_date"`
	DueDate      time.Time                      `json:"due_date"`
	Status       models.KeyResultProgressStatus `json:"status"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
	// UpdatedBy    string                         `json:"updated_by"`
}
