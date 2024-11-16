package dto

type CreateCompanyRequest struct {
	Name      string `json:"name"`
	CreatorId string `json:"creator_id"`
}