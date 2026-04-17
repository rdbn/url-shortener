package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration

	// PostgreSQL
	PostgresHost            string
	PostgresPort            int
	PostgresUser            string
	PostgresPassword        string
	PostgresDB              string
	PostgresSSLMode         string
	PostgresMaxOpenConns    int
	PostgresMaxIdleConns    int
	PostgresConnMaxLifetime time.Duration

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
	RedisTTL      time.Duration

	// App
	Environment   string
	LogLevel      string
	BaseURL       string
	FlushInterval time.Duration // интервал сброса
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
		panic(err)
	}

	return &Config{
		// Server defaults
		ServerPort:         GetEnv("SERVER_PORT", ":8080"),
		ServerReadTimeout:  GetDuration("SERVER_READ_TIMEOUT", 5*time.Second),
		ServerWriteTimeout: GetDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		ServerIdleTimeout:  GetDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),

		// PostgreSQL defaults
		PostgresHost:            GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:            GetEnvAsInt("POSTGRES_PORT_EXTERNAL", 5432),
		PostgresUser:            GetEnv("POSTGRES_USER", "user"),
		PostgresPassword:        GetEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:              GetEnv("POSTGRES_DB", "urlshortener"),
		PostgresSSLMode:         GetEnv("POSTGRES_SSLMODE", "disable"),
		PostgresMaxOpenConns:    GetEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
		PostgresMaxIdleConns:    GetEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10),
		PostgresConnMaxLifetime: GetDuration("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute),

		// Redis defaults
		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnvAsInt("REDIS_PORT", 6379),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       GetEnvAsInt("REDIS_DB", 0),
		RedisTTL:      GetDuration("REDIS_TTL", 10*time.Minute),

		// App
		Environment:   GetEnv("ENVIRONMENT", "development"),
		LogLevel:      GetEnv("LOG_LEVEL", "info"),
		BaseURL:       GetEnv("BASE_URL", "http://localhost:8080"),
		FlushInterval: 5 * time.Second,
	}
}

func (config *Config) GetPostgresDSN() string {
	hostPort := net.JoinHostPort(config.PostgresHost, strconv.Itoa(config.PostgresPort))

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.PostgresUser,
		config.PostgresPassword,
		hostPort,
		config.PostgresDB,
		config.PostgresSSLMode,
	)
}

func (config *Config) GetRedisDSN() string {
	return fmt.Sprintf("%s:%d",
		config.RedisHost,
		config.RedisPort,
	)
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intVal, err := strconv.Atoi(value)
		if err == nil {
			return intVal
		}
	}

	return defaultValue
}

func GetDuration(key string, duration time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		durationVal, err := time.ParseDuration(value)
		if err == nil {
			return durationVal
		}
	}

	return duration
}
