package http

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/dto"

	"go.uber.org/zap"
)

func (handler *URLHandler) ShortURLStats(responseWriter http.ResponseWriter, request *http.Request) {
	log := handler.logger.WithContext(request.Context())

	shortURLStats, err := dto.NewShortURLStatsInputDto(request.URL.Path)
	if err != nil {
		log.Error(err.Error())
		http.Error(responseWriter, "Not found", http.StatusNotFound)

		return
	}

	url, err := handler.URLUseCase.ShortURLStats(request.Context(), shortURLStats)
	if err != nil {
		log.Error("Not found url from storage", zap.Error(err))
		http.Error(responseWriter, "Not found url", http.StatusNotFound)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(responseWriter).Encode(url)
	if err != nil {
		log.Error("Invalidate encode json", zap.Error(err))
		http.Error(responseWriter, "Invalidate encode json", http.StatusBadRequest)

		return
	}
}
