package dto

type CreateObjectiveRequest struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Creator  string `json:"creator,omitempty"`
	Deadline string `json:"deadline,omitempty"`
}
