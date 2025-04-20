package models

import (
	"time"
)

// TODO: move to a separate file
type ShortenRequest struct {
	LongURL       string `json:"longUrl" binding:"required,url"`
	CustomShortID string `json:"customShortId,omitempty"`
}

type ShortenResponse struct {
	ShortURL  string    `json:"shortUrl"`
	LongURL   string    `json:"longUrl"`
	ShortID   string    `json:"shortId"`
	CreatedAt time.Time `json:"createdAt"`
}

type URLMapping struct {
	ShortID   string    `json:"shortId"`
	LongURL   string    `json:"longUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Hits      int64     `json:"hits"`
}
