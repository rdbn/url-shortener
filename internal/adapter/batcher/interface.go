package batcher

import (
	"context"
	"url-shortener/internal/domain"
)

type URLBatcherWorker interface {
	Start()
	AddEvent(ctx context.Context, urlDomain *domain.URL) error
	Stop()
}
