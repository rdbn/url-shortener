package redis

import (
	"context"
	"testing"
	"url-shortener/internal/domain"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func setupRedisTest(t *testing.T) *URLRedis {
	t.Helper()
	// Запускаем in-memory Redis
	miniRedis := miniredis.RunT(t)

	// Создаем реальный клиент, но подключаемся к miniredis
	client := redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})

	return &URLRedis{cache: client}
}
func TestCachedStorage_Get_CacheHit(t *testing.T) {
	cacheStore := setupRedisTest(t)
	ctx := context.Background()

	// Сохраняем данные напрямую в Redis (минуя main storage)
	url := &domain.URL{
		ShortCode:   "cachehit",
		OriginalURL: "https://example.com",
		Clicks:      42,
	}

	err := cacheStore.Save(ctx, url)
	if err != nil {
		t.Fatalf("failed to save to cache: %v", err)

		return
	}

	result, err := cacheStore.Get(ctx, "cachehit")
	if err != nil {
		t.Fatalf("error getting cache hit: %v", err)
	}

	if result.Clicks != 42 {
		t.Fatalf("wrong clicks: got %d, want 42", result.Clicks)
	}
}
