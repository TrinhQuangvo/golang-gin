package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Email        string     `json:"email" gorm:"unique;not null"`
	Password     string     `json:"-" gorm:"not null"`
	RefreshToken string     `json:"-"`
	ProfileID    *uuid.UUID `json:"-" gorm:"type:uuid"`
	Profile      *Profile   `json:"profile" gorm:"foreignKey:ProfileID;references:ID;constraint:onUpdate:CASCADE,onDelete:SET NULL;"`
	Roles        []Role     `json:"roles" gorm:"many2many:auth_roles;"`
}
