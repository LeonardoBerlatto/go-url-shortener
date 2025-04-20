package models

import (
	"time"
)

// ShortenRequest represents the input for shortening a URL
type ShortenRequest struct {
	LongURL       string `json:"longUrl" binding:"required,url"`
	CustomShortID string `json:"customShortId,omitempty"`
}

// ShortenResponse represents the API response for shortening a URL
type ShortenResponse struct {
	ShortURL  string    `json:"shortUrl"`
	LongURL   string    `json:"longUrl"`
	ShortID   string    `json:"shortId"`
	CreatedAt time.Time `json:"createdAt"`
}

// URLMapping represents the stored URL mapping
type URLMapping struct {
	ShortID   string    `json:"shortId"`
	LongURL   string    `json:"longUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Hits      int64     `json:"hits"`
}
