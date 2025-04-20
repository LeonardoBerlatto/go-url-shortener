package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leonardoberlatto/go-url-shortener/internal/models"
	"github.com/leonardoberlatto/go-url-shortener/internal/storage"
)

const (
	defaultIDLength = 8
)

var (
	// Allows only alphanumeric characters, underscores, and hyphens
	validShortIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type URLService struct {
	storage storage.Storage
	cache   *storage.RedisCache
	baseURL string
}

func NewURLService(storage storage.Storage, cache *storage.RedisCache, baseURL string) *URLService {
	return &URLService{
		storage: storage,
		cache:   cache,
		baseURL: baseURL,
	}
}

func (s *URLService) Shorten(ctx context.Context, req models.ShortenRequest) (models.ShortenResponse, error) {
	shortID := req.CustomShortID

	// If no custom ID provided, generate one
	if shortID == "" {
		generated := uuid.New().String()
		shortID = strings.ReplaceAll(generated, "-", "")[:defaultIDLength]
	} else {
		// Validate custom short ID
		if !validShortIDPattern.MatchString(shortID) {
			return models.ShortenResponse{}, fmt.Errorf("invalid custom short ID format")
		}

		// Check if the custom ID already exists
		exists, err := s.storage.CheckExists(ctx, shortID)
		if err != nil {
			return models.ShortenResponse{}, err
		}
		if exists {
			return models.ShortenResponse{}, storage.ErrorConflict
		}
	}

	// Create URL mapping
	mapping := models.URLMapping{
		ShortID:   shortID,
		LongURL:   req.LongURL,
		CreatedAt: time.Now(),
		Hits:      0,
	}

	// Store in persistent storage
	if err := s.storage.Store(ctx, mapping); err != nil {
		return models.ShortenResponse{}, err
	}

	// Store in cache
	if s.cache != nil {
		if err := s.cache.Set(ctx, mapping); err != nil {
			// Just log the error, don't fail the operation
			fmt.Printf("Error caching URL mapping: %v\n", err)
		}
	}

	return models.ShortenResponse{
		ShortURL:  fmt.Sprintf("%s/%s", s.baseURL, shortID),
		LongURL:   req.LongURL,
		ShortID:   shortID,
		CreatedAt: mapping.CreatedAt,
	}, nil
}

func (s *URLService) Resolve(ctx context.Context, shortID string) (string, error) {
	var mapping models.URLMapping
	var err error

	// Try to get from cache first
	if s.cache != nil {
		mapping, err = s.cache.Get(ctx, shortID)
		if err == nil {
			// Increment hits in background
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = s.storage.IncrementHits(ctx, shortID)
			}()
			return mapping.LongURL, nil
		}

		// Cache miss, only proceed if it's a not found error
		if err != storage.ErrorNotFound {
			fmt.Printf("Cache error: %v\n", err)
		}
	}

	// Get from storage
	mapping, err = s.storage.Get(ctx, shortID)
	if err != nil {
		return "", err
	}

	// Update cache in background
	if s.cache != nil {
		go func(m models.URLMapping) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.cache.Set(ctx, m)
		}(mapping)
	}

	go s.incrementHits(shortID)

	return mapping.LongURL, nil
}

func (s *URLService) incrementHits(shortID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.storage.IncrementHits(ctx, shortID)
}

func (s *URLService) Delete(ctx context.Context, shortID string) error {
	err := s.storage.Delete(ctx, shortID)
	if err != nil {
		return err
	}

	// Delete from cache
	if s.cache != nil {
		if err := s.cache.Delete(ctx, shortID); err != nil {
			fmt.Printf("Error deleting URL mapping from cache: %v\n", err)
		}
	}

	return nil
}
