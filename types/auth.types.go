package types

type AuthResponse struct {
	ID      uint            `json:"id"`
	Email   string          `json:"email"`
	Profile ProfileResponse `json:"profile"`
}
