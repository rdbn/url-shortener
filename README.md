# URL shortener

Небольшой HTTP-сервис для создания коротких ссылок и перенаправления на исходные URL. PostgreSQL хранит ссылки и счётчики кликов, Redis кэширует URL. При редиректе приложение атомарно увеличивает счётчик в PostgreSQL.

## Быстрый запуск

1. Создайте локальную конфигурацию: `cp .env.dist .env`.
2. Запустите PostgreSQL и Redis: `docker compose up -d postgres redis`.
3. Примените миграции: `make migrate-docker-up`.
4. Запустите сервер: `go run ./cmd/server`.

Если сервис запускается в Docker Compose, укажите в `.env` имена сервисов `postgres` и `redis` как `POSTGRES_HOST` и `REDIS_HOST`.

## HTTP API

Создать короткую ссылку:

```sh
curl -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Откройте путь с полученным кодом для перенаправления. Статистика доступна по адресу `/stats/{shortCode}`.

## Конфигурация

Приложение читает переменные окружения и необязательный `.env` файл. Список настроек и локальные значения находятся в `.env.dist`.
