package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `json:"title" gorm:"unique;not null"`
	Slug     string `json:"slug" gorm:"unique;not null"`
	Body     string `json:"body" gorm:"type:text;not null"`
	AuthorID uint   `json:"-"`
	Author   Auth   `json:"author" gorm:"foreignKey:AuthorID;references:ID"`
}
