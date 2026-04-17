package dto

import (
	"strings"
	"url-shortener/internal/domain"
)

type ShortURLStatsInputDto struct {
	ShortCode string `json:"url"`
}

func NewShortURLStatsInputDto(shortCode string) (*ShortURLStatsInputDto, error) {
	shortCode = strings.TrimPrefix(shortCode, "/stats/")

	if shortCode == "" {
		return nil, domain.ErrShortCodeNotFound
	}

	return &ShortURLStatsInputDto{ShortCode: shortCode}, nil
}
