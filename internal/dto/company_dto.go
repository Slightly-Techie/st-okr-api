package dto

type CreateCompanyRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatorId string `json:"creator_id"`
}