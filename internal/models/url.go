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

type PaginationRequest struct {
	PageNumber int `form:"pageNumber,default=1" binding:"min=1"`
	PageSize   int `form:"pageSize,default=10" binding:"min=1,max=100"`
}

type PaginatedURLsResponse struct {
	Content    []URLMapping `json:"content"`
	TotalCount int64        `json:"totalCount"`
	PageNumber int          `json:"pageNumber"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}
