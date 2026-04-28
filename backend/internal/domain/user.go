package domain

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	ConnpassID string `json:"connpass_id,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
