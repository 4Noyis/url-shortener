package models

import "time"

type URL struct {
	ID        int       `json:"id"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Clicks    int       `json:"clicks"`
	UserID    *string   `json:"user_id,omitempty"`
	IsActive  bool      `json:"is_active"`
}
