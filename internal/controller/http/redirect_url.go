package http

import (
	"errors"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"

	"go.uber.org/zap"
)

func (handler *URLHandler) RedirectURL(responseWriter http.ResponseWriter, request *http.Request) {
	log := handler.logger.WithContext(request.Context())

	redirectURLInputDto, err := dto.NewRedirectURLInputDto(request.URL.Path)
	if err != nil {
		log.Error(err.Error())
		http.Error(responseWriter, "Not found", http.StatusNotFound)

		return
	}

	redirectURL, err := handler.URLUseCase.RedirectURL(request.Context(), redirectURLInputDto)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Error("Not found url from storage", zap.Error(err))
			http.Error(responseWriter, "Not found", http.StatusNotFound)

			return
		}

		log.Error("Failed to increment clicks", zap.Error(err))
		http.Error(responseWriter, "Failed to increment clicks", http.StatusBadGateway)

		return
	}

	log.Info(
		"Redirecting to original URL",
		zap.String("short_code", redirectURL.ShortCode),
		zap.String("original_url", redirectURL.OriginalURL),
	)

	http.Redirect(responseWriter, request, redirectURL.OriginalURL, http.StatusFound)
}
