# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Устанавливаем необходимые инструменты
RUN apk add --no-cache git make

# Устанавливаем golangci-lint для проверки кода при сборке
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.55.2

WORKDIR /var/www/app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

# Копируем исходный код
COPY . .

# Собираем приложение
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/url-shortener ./cmd/server

# --- ЭТАП 2: Для разработки (Development) ---
FROM builder AS dev

# Исправленная команда установки
RUN go install github.com/air-verse/air@latest

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

USER appuser

# Запуск
CMD ["air", "-c", ".air.toml"]
