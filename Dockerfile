# Стадия сборки
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git

# Копируем сначала только модули для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем swag и swagger-ui
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN apk add --no-cache curl && \
    mkdir -p /swagger-ui && \
    curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.9.0.tar.gz | tar xz -C /swagger-ui --strip-components=1 swagger-ui-5.9.0/dist

# Копируем остальные файлы
COPY . .

# Генерируем Swagger документацию
RUN swag init -g ./cmd/main.go --output ./docs --parseDependency --parseInternal --parseDepth 2

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /popcorntime ./cmd/

# Финальная стадия
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и необходимые файлы
COPY --from=builder /popcorntime .
COPY --from=builder /app/docs ./docs
COPY --from=builder /swagger-ui ./swagger-ui

# Настраиваем права
RUN chmod -R 755 /app/docs /app/swagger-ui

# Остальные настройки остаются без изменений
...