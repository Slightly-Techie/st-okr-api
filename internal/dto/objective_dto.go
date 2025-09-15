package dto

import (
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
)

type CreateObjectiveRequest struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description"`
	Type        models.ObjectiveType  `json:"type" validate:"required,oneof=company team"`
	OwnerID     string                `json:"owner_id" validate:"required,uuid"`
	CompanyID   string                `json:"company_id" validate:"required,uuid"`
	TeamID      *string               `json:"team_id,omitempty" validate:"omitempty,uuid"`
	StartDate   time.Time             `json:"start_date" validate:"required"`
	EndDate     time.Time             `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateObjectiveRequest struct {
	ID          string                `json:"id" validate:"required,uuid"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Status      models.ObjectiveStatus `json:"status" validate:"omitempty,oneof=draft active completed archived on_hold"`
	StartDate   time.Time             `json:"start_date"`
	EndDate     time.Time             `json:"end_date"`
}

type ObjectiveResponse struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Type        models.ObjectiveType   `json:"type"`
	OwnerID     string                 `json:"owner_id"`
	CompanyID   string                 `json:"company_id"`
	TeamID      *string                `json:"team_id,omitempty"`
	Status      models.ObjectiveStatus `json:"status"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     time.Time              `json:"end_date"`
	Progress    float64                `json:"progress"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	KeyResults  []KeyResultResponse    `json:"key_results,omitempty"`
}

type ObjectiveListResponse struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Type        models.ObjectiveType   `json:"type"`
	Status      models.ObjectiveStatus `json:"status"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     time.Time              `json:"end_date"`
	Progress    float64                `json:"progress"`
	KeyResultsCount int                `json:"key_results_count"`
}