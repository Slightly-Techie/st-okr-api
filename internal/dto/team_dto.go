package dto

type CreateTeam struct {
	Name        string `json:"name" validate:"required,uuid"`
	CompanyID   string `json:"company_id" validate:"required,uuid"`
	Description string `json:"description"`
}

type UpdateTeam struct {
	ID          string `json:"id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type AddTeamMember struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	TeamID string `json:"team_id" validate:"required,uuid"`
}
