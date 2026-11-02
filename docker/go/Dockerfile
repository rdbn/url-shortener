# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Устанавливаем необходимые инструменты
RUN apk add --no-cache git make

# Устанавливаем golangci-lint для команды make lint
RUN wget -O- -nv https://golangci-lint.run/install.sh | sh -s v2.15.0
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /var/www/app

# Меняем владельца рабочей директории
RUN chown -R appuser:appgroup /var/www/app
RUN chown -R appuser:appgroup /go

USER appuser

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

# Копируем исходный код
COPY .. .

# --- ЭТАП 2: Для разработки (Development) ---
FROM builder AS dev

# Исправленная команда установки
RUN go install github.com/air-verse/air@latest

# Запуск
CMD ["air", "-c", ".air.toml"]
