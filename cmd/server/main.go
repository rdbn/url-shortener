package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/config"
	"url-shortener/internal/adapter/batcher"
	"url-shortener/internal/adapter/postgres"
	"url-shortener/internal/adapter/redis"
	http2 "url-shortener/internal/controller/http"
	"url-shortener/internal/usecase"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Environment, cfg.LogLevel)
	defer log.Sync()

	log.Info("Starting server",
		zap.String("environment", cfg.Environment),
		zap.String("log_level", cfg.LogLevel),
		zap.String("port", cfg.ServerPort),
	)

	urlPostgres := postgres.NewURLPostgres(
		cfg.GetPostgresDSN(),
		cfg.PostgresMaxOpenConns,
		cfg.PostgresMaxIdleConns,
		cfg.PostgresConnMaxLifetime,
	)
	defer urlPostgres.Close()

	urlRedis := redis.NewURLRedis(
		cfg.GetRedisDSN(),
		cfg.RedisPassword,
		cfg.RedisDB,
		cfg.RedisTTL,
	)
	defer urlRedis.Close()

	urlBatcher := batcher.NewURLBatcher(
		urlPostgres,
		urlRedis,
		log,
		cfg.FlushInterval,
	)
	defer urlBatcher.Stop()

	urlBatcher.Start()
	urlUseCase := usecase.NewURLUseCase(urlPostgres, urlRedis, urlBatcher)
	mux := http2.Router(urlUseCase, log, cfg.BaseURL)

	recoveryMiddleware := middleware.RecoveryMiddleware(log)(
		middleware.LoggerMiddleware(log)(
			middleware.CORSMiddleware([]string{"*"})(mux),
		),
	)

	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      recoveryMiddleware,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	log.Info("Server starting on: " + cfg.ServerPort)

	err := server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Server failed:", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error("Server shutdown failed:", zap.Error(err))
	}
}
