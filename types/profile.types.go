package types

type ProfileResponse struct {
	ID           uint   `json:"id"`
	FriendlyName string `json:"friendly_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	Avartar      string `json:"avatar"`
}
