package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/leonardoberlatto/go-url-shortener/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL = time.Hour
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, shortID string) (models.URLMapping, error) {
	data, err := r.client.Get(ctx, shortID).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.URLMapping{}, ErrorNotFound
		}
		return models.URLMapping{}, err
	}

	var mapping models.URLMapping
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return models.URLMapping{}, err
	}

	return mapping, nil
}

func (r *RedisCache) Set(ctx context.Context, mapping models.URLMapping) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, mapping.ShortID, data, cacheTTL).Err()
}

func (r *RedisCache) Delete(ctx context.Context, shortID string) error {
	return r.client.Del(ctx, shortID).Err()
}

func (r *RedisCache) CheckExists(ctx context.Context, shortID string) (bool, error) {
	exists, err := r.client.Exists(ctx, shortID).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
