package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"url-shortener/internal/domain"

	"github.com/go-redis/redis/v8"
)

const cacheKey = "url_domain_batch"

type URLRedis struct {
	cache *redis.Client
	ttl   time.Duration
}

func NewURLRedis(redisDSN string, password string, db int, ttl time.Duration) *URLRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisDSN,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return &URLRedis{cache: client, ttl: ttl}
}

func (c *URLRedis) Get(ctx context.Context, shortCode string) (*domain.URL, error) {
	data, err := c.cache.Get(ctx, c.cacheKey(shortCode)).Result()

	if errors.Is(err, redis.Nil) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	var url domain.URL

	err = json.Unmarshal([]byte(data), &url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (c *URLRedis) Save(ctx context.Context, url *domain.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	err = c.cache.Set(ctx, c.cacheKey(url.ShortCode), data, c.ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *URLRedis) Delete(ctx context.Context, shortCode string) error {
	err := c.cache.Del(ctx, c.cacheKey(shortCode)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *URLRedis) GetURLDomainBatchAndClear(ctx context.Context) (map[string]domain.URL, error) {
	pipe := c.cache.Pipeline()
	getAll := pipe.HGetAll(ctx, cacheKey)
	pipe.Del(ctx, cacheKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	batch, err := getAll.Result()
	if err != nil {
		return nil, err
	}

	clickBatch := make(map[string]domain.URL)

	for _, v := range batch {
		var url domain.URL

		err := json.Unmarshal([]byte(v), &url)
		if err != nil {
			return nil, err
		}

		clickBatch[url.ShortCode] = url
	}

	return clickBatch, nil
}

func (c *URLRedis) SaveToBatch(ctx context.Context, url *domain.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	err = c.cache.HSet(ctx, cacheKey, c.cacheKey(url.ShortCode), data).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *URLRedis) Close() {
	err := c.cache.Close()
	if err != nil {
		panic(err)
	}
}

func (c *URLRedis) cacheKey(shortCode string) string {
	return "url:" + shortCode
}
