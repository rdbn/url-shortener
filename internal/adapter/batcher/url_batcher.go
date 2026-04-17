package batcher

import (
	"context"
	"sync"
	"time"
	"url-shortener/internal/adapter/postgres"
	"url-shortener/internal/adapter/redis"
	"url-shortener/internal/domain"
	"url-shortener/pkg/logger"

	"go.uber.org/zap"
)

type URLBatcher struct {
	postgresStorage postgres.Storage
	redisStorage    redis.Storage
	logger          *logger.Logger
	flushInterval   time.Duration
	stopCh          chan struct{}
	wg              *sync.WaitGroup
}

func NewURLBatcher(
	postgresStorage postgres.Storage,
	redisStorage redis.Storage,
	logger *logger.Logger,
	flushInterval time.Duration,
) *URLBatcher {
	b := &URLBatcher{
		postgresStorage: postgresStorage,
		redisStorage:    redisStorage,
		logger:          logger,
		flushInterval:   flushInterval,
		stopCh:          make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}

	return b
}

func (b *URLBatcher) Start() {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		ticker := time.NewTicker(b.flushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				b.saveBatch(ctx)
				cancel()
			case <-b.stopCh:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				b.saveBatch(ctx)
				cancel()
			}
		}
	}()
}

func (b *URLBatcher) AddEvent(ctx context.Context, urlDomain *domain.URL) error {
	err := b.redisStorage.SaveToBatch(ctx, urlDomain)
	if err != nil {
		return err
	}

	return nil
}

func (b *URLBatcher) Stop() {
	b.logger.Info("stopping click batcher...")

	close(b.stopCh)
	b.wg.Wait()
}

func (b *URLBatcher) saveBatch(ctx context.Context) {
	urlDomainBatch, err := b.redisStorage.GetURLDomainBatchAndClear(ctx)
	if err != nil {
		b.logger.Error("failed to create click batch", zap.Error(err))

		return
	}

	if len(urlDomainBatch) == 0 {
		return
	}

	start := time.Now()

	err = b.postgresStorage.SaveBatch(ctx, urlDomainBatch)
	if err != nil {
		b.logger.Error("failed to save batch", zap.Error(err), zap.Int("batch_size", len(urlDomainBatch)))

		return
	}

	b.logger.Debug("batch saved",
		zap.Int("batch_size", len(urlDomainBatch)),
		zap.Duration("duration", time.Since(start)),
	)
}
