package dto

import (
	"encoding/json"
	"io"
	"url-shortener/internal/domain"
)

type CreateURLInputDto struct {
	URL string `json:"url"`
}

func NewCreateURLInputDto(requestBody io.ReadCloser) (*CreateURLInputDto, error) {
	createURLInputDto := &CreateURLInputDto{}

	err := json.NewDecoder(requestBody).Decode(&createURLInputDto)
	if err != nil {
		return nil, err
	}

	if createURLInputDto.URL == "" {
		return nil, domain.ErrURLIsRequired
	}

	return createURLInputDto, nil
}

type CreateURLOutputDto struct {
	ShortURL  string `json:"shortUrl"`
	ShortCode string `json:"shortCode"`
}

func NewCreateURLOutputDto(shortURL, shortCode string) *CreateURLOutputDto {
	return &CreateURLOutputDto{
		ShortURL:  shortURL,
		ShortCode: shortCode,
	}
}
