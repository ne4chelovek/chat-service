FROM golang:1.24-bullseye AS builder

# Установка зависимостей для Kafka (C-библиотеки)
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

# Скачивание Go-модулей
RUN go mod download

# Сборка приложения (CGO обязателен для confluent-kafka-go)
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/mikle-chat cmd/main.go

# Финальный образ (минимальный)
FROM debian:bullseye-slim

# Установка runtime-зависимостей для Kafka
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

# Создание папки для сертификатов и копирование
WORKDIR /root/
RUN mkdir -p certs
COPY --from=builder /app/bin/mikle-chat .
COPY certs/. ./certs/

# Запуск приложения
CMD ["./mikle-chat"]