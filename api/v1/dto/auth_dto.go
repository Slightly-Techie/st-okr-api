package dto

type AuthResponse struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
