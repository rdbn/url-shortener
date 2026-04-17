package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/dto"

	"go.uber.org/zap"
)

func (handler *URLHandler) CreateURL(responseWriter http.ResponseWriter, request *http.Request) {
	log := handler.logger.WithContext(request.Context())

	if request.Method != http.MethodPost {
		log.Warn("method not allowed", zap.String("method", request.Method))
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	createURLInputDto, err := dto.NewCreateURLInputDto(request.Body)
	if err != nil {
		if errors.Is(err, domain.ErrURLIsRequired) {
			http.Error(responseWriter, "Url is required", http.StatusBadRequest)

			return
		}

		log.Error("failed to decode request", zap.Error(err))
		http.Error(responseWriter, "Invalidate request body", http.StatusBadRequest)

		return
	}

	url, err := handler.URLUseCase.CreateURL(request.Context(), createURLInputDto)
	if err != nil {
		log.Error("failed to create URL", zap.Error(err))
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)

		return
	}

	log.Info("URL created",
		zap.String("short_code", url.ShortCode),
		zap.String("original_url", url.OriginalURL),
	)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	createURLOutputDto := dto.NewCreateURLOutputDto(handler.baseURL+"/"+url.ShortCode, url.ShortCode)

	err = json.NewEncoder(responseWriter).Encode(createURLOutputDto)
	if err != nil {
		log.Error("Invalidate encode json", zap.Error(err))
		http.Error(responseWriter, "Invalidate encode json", http.StatusBadRequest)

		return
	}
}
