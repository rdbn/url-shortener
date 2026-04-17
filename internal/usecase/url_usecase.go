package usecase

import (
	"url-shortener/internal/adapter/batcher"
	"url-shortener/internal/adapter/postgres"
	"url-shortener/internal/adapter/redis"
)

type URLUseCase struct {
	urlRedis    redis.Storage
	urlPostgres postgres.Storage
	batcher     *batcher.URLBatcher
}

func NewURLUseCase(urlPostgres postgres.Storage, urlRedis redis.Storage, batcher *batcher.URLBatcher) *URLUseCase {
	return &URLUseCase{
		urlRedis:    urlRedis,
		urlPostgres: urlPostgres,
		batcher:     batcher,
	}
}
