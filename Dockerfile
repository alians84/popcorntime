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
WORKDIR /app
RUN swag init -g ./cmd/main.go --output ./docs --parseDependency --parseInternal --parseDepth 2

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /popcorntime ./cmd/

# Финальная стадия
FROM alpine:latest

# Копируем бинарник и необходимые файлы
COPY --from=builder /popcorntime .
COPY --from=builder /app/docs ./docs
COPY --from=builder /swagger-ui ./swagger-ui

# Настраиваем права
RUN chmod -R 755 /app/docs /app/swagger-ui
# Устанавливаем зависимости
RUN apk add --no-cache tzdata ca-certificates

# Добавляем переменную окружения для JWT
ENV JWT_SECRET=01f8b938cea0d104a928348ce7682f0f4b693d88faff8ba70d557e3fdb6a5f87af1a0531f21add905a6fbb1577db6231a53d54ab0bece1db2595f6bfbc862a3a2bbb0d5f4b75ade1c839dea36dc877be0b39f15a93a04f544f6296def70f1b2588d2d9560ddca3bdc33ad226a645cf18bf4657484b37191febd4ca9fd1c4a31e242dfa6f57c106532cfdd7b88b0a1bfd6d92e5b96be6256e04202efe8c7ff9354e75326fa9f8db52a686259da56365bb2bd0f0cfcc27415c7f786e4c212de0d32bde7af04d4581660b53ddb250689714c514b7b8fb4832360149cf9597b19b8391ceb35f5196f29ca0b45b956c803edf9e4c4d9b481ebbe8239cec50ef9dd857

EXPOSE 3000

CMD ["/app/popcorntime"]