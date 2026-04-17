package http

import (
	"net/http"
	"url-shortener/internal/usecase"
	"url-shortener/pkg/logger"
)

func Router(urlUseCase *usecase.URLUseCase, logger *logger.Logger, baseURL string) http.Handler {
	urlHandler := NewURLHandler(urlUseCase, logger, baseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", urlHandler.CreateURL)
	mux.HandleFunc("/stats/", urlHandler.ShortURLStats)
	mux.HandleFunc("/", urlHandler.RedirectURL)

	return mux
}
