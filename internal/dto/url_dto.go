package dto

import "time"

type ShortenURLRequest struct {
	LongURL    string `json:"long_url" binding:"required,url"`
	TTLSeconds *int   `json:"ttl_seconds,omitempty" binding:"omitempty,min=60,max=31536000"` // 1 minute to 1 year
}

type ShortenURLResponse struct {
	ShortURL  string     `json:"short_url"`
	LongURL   string     `json:"long_url"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}