package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	FriendlyName string `json:"friendly_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	Avartar      string `json:"avatar"`
	Status       bool   `json:"status"`
}
