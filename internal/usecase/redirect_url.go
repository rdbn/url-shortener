package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"

	"golang.org/x/sync/singleflight"
)

func (useCase *URLUseCase) RedirectURL(ctx context.Context, dto *dto.RedirectURLInputDto) (*domain.URL, error) {
	sf := &singleflight.Group{}

	val, err, _ := sf.Do(dto.ShortCode, func() (any, error) {
		urlDomain, err := useCase.urlRedis.Get(ctx, dto.ShortCode)
		if errors.Is(err, domain.ErrNotFound) {
			urlDomain, err = useCase.urlPostgres.Get(ctx, dto.ShortCode)

			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrNotFound
			}

			if err != nil {
				return nil, err
			}
		}

		if err != nil {
			return nil, err
		}

		urlDomain.Clicks++
		urlDomain.UpdatedAt = time.Now()

		err = useCase.urlRedis.Save(ctx, urlDomain)
		if err != nil {
			return nil, err
		}

		err = useCase.batcher.AddEvent(ctx, urlDomain)
		if err != nil {
			return nil, err
		}

		return urlDomain, nil
	})
	if err != nil {
		return nil, err
	}

	urlDomain, ok := val.(*domain.URL)
	if !ok {
		return nil, fmt.Errorf("%w: %T", domain.ErrUnexpectedType, val)
	}

	return urlDomain, nil
}
