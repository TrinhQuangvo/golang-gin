package types

import (
	"time"
)

type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	Total       int         `json:"total"`
	Total_Pages int         `json:"total_pages"`
}

type PostResponse struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Slug      string       `json:"slug"`
	Body      string       `json:"body"`
	AuthID    uint         `json:"auth_id"`
	Auth      AuthResponse `json:"auth"`
	CreatedAt time.Time    `json:"created_at"`
}
