package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"
	"url-shortener/internal/utils"
)

func (useCase *URLUseCase) CreateURL(ctx context.Context, dto *dto.CreateURLInputDto) (*domain.URL, error) {
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	shortCode := createShortURL(dto.URL)

	urlDomain := &domain.URL{
		ID:          uuid,
		ShortCode:   shortCode,
		OriginalURL: dto.URL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Clicks:      0,
	}

	err = useCase.urlRedis.Save(ctx, urlDomain)
	if err != nil {
		return nil, err
	}

	err = useCase.urlPostgres.Save(ctx, urlDomain)
	if err != nil {
		errRedisDelete := useCase.urlRedis.Delete(ctx, urlDomain.ShortCode)
		if errRedisDelete != nil {
			return nil, errRedisDelete
		}
		return nil, err
	}

	return urlDomain, nil
}

func createShortURL(originalURL string) string {
	data := []byte(originalURL + time.Now().String())
	hash := sha256.Sum256(data)

	code := base64.URLEncoding.EncodeToString(hash[:6])

	return code[:6]
}
