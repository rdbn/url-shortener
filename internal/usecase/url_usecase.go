package usecase

import (
	"url-shortener/internal/adapter/postgres"
	"url-shortener/internal/adapter/redis"
)

type URLUseCase struct {
	urlRedis    redis.Storage
	urlPostgres postgres.Storage
}

func NewURLUseCase(urlPostgres postgres.Storage, urlRedis redis.Storage) *URLUseCase {
	return &URLUseCase{
		urlRedis:    urlRedis,
		urlPostgres: urlPostgres,
	}
}
