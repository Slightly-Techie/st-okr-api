package dto

import (
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
)

type CreateKeyResultRequest struct {
	ObjectiveID  string              `json:"objective_id" validate:"required,uuid"`
	Title        string              `json:"title" validate:"required"`
	Description  string              `json:"description"`
	MetricType   models.MetricType   `json:"metric_type" validate:"required,metric_type"`
	CurrentValue float64             `json:"current_value"`
	TargetValue  float64             `json:"target_value" validate:"required"`
	AssigneeType models.AssigneeType `json:"assignee_type" validate:"required,assignee_type"`
	AssigneeID   string              `json:"assignee_id" validate:"required,uuid"`
	StartDate    time.Time           `json:"start_date" validate:"required"`
	DueDate      time.Time           `json:"due_date" validate:"due_date"`
}

type UpdateKeyResultRequest struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	CurrentValue float64             `json:"current_value"`
	TargetValue  float64             `json:"target_value"`
	MetricType   models.MetricType   `json:"metric_type"`
	AssigneeType models.AssigneeType `json:"assignee_type" validate:"assignee_type"`
	AssigneeID   string              `json:"assignee_id" validate:"uuid"`
	StartDate    time.Time           `json:"start_date"`
	DueDate      time.Time           `json:"due_date" validate:"due_date"`
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
}
