package storage

import (
	"context"

	"github.com/leonardoberlatto/go-url-shortener/internal/models"
)

type Storage interface {
	Store(ctx context.Context, mapping models.URLMapping) error

	Get(ctx context.Context, shortID string) (models.URLMapping, error)

	Delete(ctx context.Context, shortID string) error

	CheckExists(ctx context.Context, shortID string) (bool, error)

	IncrementHits(ctx context.Context, shortID string) error

	ListURLs(ctx context.Context, pageNumber, pageSize int) ([]models.URLMapping, int64, error)
}
