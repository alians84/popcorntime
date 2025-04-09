# Билд стадии
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости и swag для документации
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем исходники
COPY . .

# Генерируем Swagger-документацию
RUN swag init -g cmd/main.go --output docs  --parseDependency --parseInternal --parseDepth 2

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /popcorntime

# Финальная стадия
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и конфиги
COPY --from=builder /popcorntime /app/
COPY --from=builder /app/config /app/config

# Устанавливаем tzdata для работы с временными зонами
RUN apk add --no-cache tzdata

EXPOSE 3000

CMD ["/app/popcorntime"]