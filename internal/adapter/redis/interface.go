package redis

import (
	"context"
	"url-shortener/internal/domain"
)

type Storage interface {
	Get(ctx context.Context, shortCode string) (*domain.URL, error)
	Save(ctx context.Context, url *domain.URL) error
	Delete(ctx context.Context, shortCode string) error
	GetURLDomainBatchAndClear(ctx context.Context) (map[string]domain.URL, error)
	SaveToBatch(ctx context.Context, url *domain.URL) error
	Close()
}
