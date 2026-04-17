package http

import (
	"url-shortener/internal/usecase"
	"url-shortener/pkg/logger"
)

type URLHandler struct {
	URLUseCase *usecase.URLUseCase
	logger     *logger.Logger
	baseURL    string
}

func NewURLHandler(urlUseCase *usecase.URLUseCase, logger *logger.Logger, baseURL string) *URLHandler {
	return &URLHandler{URLUseCase: urlUseCase, logger: logger, baseURL: baseURL}
}
