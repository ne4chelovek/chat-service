FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/mikle-chat cmd/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/
RUN mkdir -p certs
COPY --from=builder /app/bin/mikle-chat .
COPY certs/ ./certs/

CMD ["./mikle-chat"]