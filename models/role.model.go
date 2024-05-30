package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
	Slug string `json:"slug" gorm:"unique"`
}
