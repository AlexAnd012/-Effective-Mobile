# этап сборки
FROM golang:1.24.7-alpine AS builder
# Рабочая папка внутри контейнера, куда будем копировать код
WORKDIR /src

# кеш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# собираем статичный бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o /out/subs-api ./cmd/main.go

# рантайм
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -g '' app
WORKDIR /app

COPY --from=builder /out/subs-api /usr/local/bin/subs-api

EXPOSE 8080
USER app
ENTRYPOINT ["/usr/local/bin/subs-api"]
