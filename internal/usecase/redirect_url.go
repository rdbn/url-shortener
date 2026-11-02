package usecase

import (
	"context"
	"errors"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"
)

func (useCase *URLUseCase) RedirectURL(ctx context.Context, dto *dto.RedirectURLInputDto) (*domain.URL, error) {
	urlDomain, err := useCase.urlRedis.Get(ctx, dto.ShortCode)
	if errors.Is(err, domain.ErrNotFound) {
		urlDomain, err = useCase.urlPostgres.Get(ctx, dto.ShortCode)
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		if err == nil {
			if err := useCase.urlRedis.Save(ctx, urlDomain); err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}

	clicks, err := useCase.urlPostgres.IncrementClicks(ctx, dto.ShortCode)
	if err != nil {
		return nil, err
	}
	urlDomain.Clicks = clicks
	return urlDomain, nil
}
