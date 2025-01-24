package dto

type CreateObjectivesRequest struct {
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	CreatorID string `json:"creator_id,omitempty"`
	Deadline  string `json:"deadline,omitempty"`
}
