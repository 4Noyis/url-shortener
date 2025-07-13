package dto

import "time"

type ShortenURLRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

type ShortenURLResponse struct {
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}