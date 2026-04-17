package postgres

import (
	"context"
	"url-shortener/internal/domain"
)

type Storage interface {
	Save(ctx context.Context, url *domain.URL) error
	Get(ctx context.Context, shortCode string) (*domain.URL, error)
	SaveBatch(ctx context.Context, batch map[string]domain.URL) error
	Close()
}
