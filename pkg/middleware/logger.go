package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
	"url-shortener/pkg/logger"

	"go.uber.org/zap"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
)

func LoggerMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			startNow := time.Now()

			requestID := GenerateRequestID()
			ctx := context.WithValue(request.Context(), RequestIDKey, requestID)
			request = request.WithContext(ctx)

			responseWriter.Header().Set("X-Request-Id", requestID)

			newResponseWriter := NewResponseWriter(responseWriter)
			next.ServeHTTP(newResponseWriter, request)

			log.WithContext(ctx).LogHTTPRequest(logger.RequestLog{
				Method:     request.Method,
				Path:       request.URL.Path,
				StatusCode: newResponseWriter.StatusCode(),
				Duration:   time.Since(startNow),
				ClientIP:   getClientIP(request),
				UserAgent:  request.UserAgent(),
				RequestID:  requestID,
			})
		})
	}
}

func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func(request *http.Request) {
				if err := recover(); err != nil {
					log.WithContext(request.Context()).Error(
						"panic recovered",
						zap.String("method", request.Method),
						zap.String("path", request.URL.Path),
					)
					http.Error(
						responseWriter,
						"Internal Server Error",
						http.StatusInternalServerError,
					)
				}
			}(request)

			next.ServeHTTP(responseWriter, request)
		})
	}
}

// CORSMiddleware добавляет CORS заголовки.
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			origin := request.Header.Get("Origin")

			// Проверяем разрешен ли origin
			allowed := false

			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true

					break
				}
			}

			if allowed {
				responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
				responseWriter.Header().
					Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				responseWriter.Header().
					Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				responseWriter.Header().Set("Access-Control-Max-Age", "86400")
			}

			if request.Method == http.MethodOptions {
				responseWriter.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(responseWriter, request)
		})
	}
}

func GenerateRequestID() string {
	bytes := make([]byte, 8)

	_, err := rand.Read(bytes)
	if err != nil {
		return hex.EncodeToString([]byte(time.Now().String()))
	}

	return hex.EncodeToString(bytes)
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	return r.RemoteAddr
}

// ResponseWriter обертка для захвата статус кода.
type ResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}
