package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/config"
	"url-shortener/internal/adapter/postgres"
	"url-shortener/internal/adapter/redis"
	http2 "url-shortener/internal/controller/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/usecase"
	"url-shortener/pkg/logger"
)

type createResultStruct struct {
	ShortURL  string `json:"shortUrl"`
	ShortCode string `json:"shortCode"`
}

func TestCreateURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full flow test in short mode")
	}

	runServer(t, func(t *testing.T, server *httptest.Server) {
		t.Helper()

		client := &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		createURL(t, client, server)
	})
}

func TestClicks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full flow test in short mode")
	}

	runServer(t, func(t *testing.T, server *httptest.Server) {
		t.Helper()

		client := &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// createURL
		createResult := createURL(t, client, server)

		// clicks
		statusResp, err := client.Get(server.URL + "/stats/" + createResult.ShortCode)
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		var stats domain.URL

		err = json.NewDecoder(statusResp.Body).Decode(&stats)
		if err != nil {
			t.Fatalf("Failed decode response create: %v", err)
		}

		if stats.Clicks != 0 {
			t.Fatalf("Expected clicks 0, got %d", stats.Clicks)
		}

		err = statusResp.Body.Close()
		if err != nil {
			t.Fatalf("Failed create url short: %v", err)
		}

		redirectResp, err := client.Get(server.URL + "/" + createResult.ShortCode)
		if err != nil {
			t.Fatalf("Failed to follow short URL: %v", err)
		}
		if redirectResp.StatusCode != http.StatusFound {
			t.Fatalf("Expected redirect status %d, got %d", http.StatusFound, redirectResp.StatusCode)
		}
		_ = redirectResp.Body.Close()

		statusResp, err = client.Get(server.URL + "/stats/" + createResult.ShortCode)
		if err != nil {
			t.Fatalf("Failed to get updated stats: %v", err)
		}
		if err := json.NewDecoder(statusResp.Body).Decode(&stats); err != nil {
			t.Fatalf("Failed to decode updated stats: %v", err)
		}
		_ = statusResp.Body.Close()
		if stats.Clicks != 1 {
			t.Fatalf("Expected clicks 1 after one redirect, got %d", stats.Clicks)
		}
	})
}

func TestRedirectURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full flow test in short mode")
	}

	runServer(t, func(t *testing.T, server *httptest.Server) {
		t.Helper()

		client := &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// createURL
		createResult := createURL(t, client, server)

		// Redirect to original URL
		redirectResp, err := client.Get(server.URL + "/" + createResult.ShortCode)
		if err != nil {
			t.Fatalf("Failed to redirect: %v", err)
		}

		err = redirectResp.Body.Close()
		if err != nil {
			t.Fatalf("Failed create url short: %v", err)
		}

		if redirectResp.StatusCode != http.StatusFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusFound, redirectResp.StatusCode)
		}

		if redirectResp.Header.Get("Location") != "https://example.com" {
			t.Fatalf(
				"Expected location %s, got %s",
				"https://example.com",
				redirectResp.Header.Get("Location"),
			)
		}
	})
}

func createURL(t *testing.T, client *http.Client, server *httptest.Server) *createResultStruct {
	t.Helper()

	createResp, err := client.Post(
		server.URL+"/shorten",
		"application/json",
		bytes.NewBufferString(`{"url": "https://example.com"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, createResp.StatusCode)
	}

	createResult := &createResultStruct{}

	err = json.NewDecoder(createResp.Body).Decode(&createResult)
	if err != nil {
		t.Fatalf("Failed decode response create: %v", err)
	}

	if createResult.ShortURL == "" {
		t.Fatal("Expected short url not empty")
	}

	err = createResp.Body.Close()
	if err != nil {
		t.Fatalf("Failed create url short: %v", err)
	}

	return createResult
}

func testConfig() *config.Config {
	return &config.Config{
		// Server defaults
		ServerPort:         config.GetEnv("SERVER_PORT", ":8080"),
		ServerReadTimeout:  config.GetDuration("SERVER_READ_TIMEOUT", 5*time.Second),
		ServerWriteTimeout: config.GetDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		ServerIdleTimeout:  config.GetDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),

		// PostgreSQL defaults
		PostgresHost:            config.GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:            config.GetEnvAsInt("POSTGRES_PORT_EXTERNAL_TEST", 5436),
		PostgresUser:            config.GetEnv("POSTGRES_USER", "urlshortener"),
		PostgresPassword:        config.GetEnv("POSTGRES_PASSWORD", "urlshortener"),
		PostgresDB:              config.GetEnv("POSTGRES_DB", "urlshortener"),
		PostgresSSLMode:         config.GetEnv("POSTGRES_SSLMODE", "disable"),
		PostgresMaxOpenConns:    config.GetEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
		PostgresMaxIdleConns:    config.GetEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10),
		PostgresConnMaxLifetime: config.GetDuration("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute),

		// Redis defaults
		RedisHost:     config.GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     config.GetEnvAsInt("REDIS_PORT_TEST", 6382),
		RedisPassword: config.GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       config.GetEnvAsInt("REDIS_DB", 0),
		RedisTTL:      config.GetDuration("REDIS_TTL", 10*time.Minute),

		BaseURL: config.GetEnv("BASE_URL", "http://localhost:8080"),
	}
}

func runServer(t *testing.T, testFunc func(*testing.T, *httptest.Server)) {
	t.Helper()

	cfg := testConfig()

	nop := logger.NewNop()

	urlPostgresStorage := postgres.NewURLPostgres(
		cfg.GetPostgresDSN(),
		cfg.PostgresMaxOpenConns,
		cfg.PostgresMaxIdleConns,
		cfg.PostgresConnMaxLifetime,
	)
	defer urlPostgresStorage.Close()

	urlRedisStorage := redis.NewURLRedis(
		cfg.GetRedisDSN(),
		cfg.RedisPassword,
		cfg.RedisDB,
		cfg.RedisTTL,
	)
	defer urlRedisStorage.Close()

	urlUseCase := usecase.NewURLUseCase(urlPostgresStorage, urlRedisStorage)
	mux := http2.Router(urlUseCase, nop, cfg.BaseURL)

	server := httptest.NewServer(mux)
	defer server.Close()

	testFunc(t, server)
}
