package usecase

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"
)

func (useCase *URLUseCase) ShortURLStats(
	ctx context.Context,
	shortURLStats *dto.ShortURLStatsInputDto,
) (*domain.URL, error) {
	url, err := useCase.urlPostgres.Get(ctx, shortURLStats.ShortCode)
	if err != nil {
		return nil, err
	}

	return url, nil
}
