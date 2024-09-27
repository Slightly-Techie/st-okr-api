package dto

type CreateCompanyRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"company_code"`
	CreatorId string `json:"creator_id"`
}
