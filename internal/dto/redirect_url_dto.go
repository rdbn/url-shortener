package dto

import (
	"strings"
	"url-shortener/internal/domain"
)

type RedirectURLInputDto struct {
	ShortCode string `json:"url"`
}

func NewRedirectURLInputDto(shortCode string) (*RedirectURLInputDto, error) {
	shortCode = strings.TrimPrefix(shortCode, "/")

	if shortCode == "" {
		return nil, domain.ErrShortCodeNotFound
	}

	return &RedirectURLInputDto{ShortCode: shortCode}, nil
}
