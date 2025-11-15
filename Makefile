APP      := subs-api  # Имя бинарника и Docker-образа
MAIN     := ./cmd/main.go
IMG      := $(APP):latest
BIN_DIR  := bin # Кладем бинарник
DB_DSN  ?= postgres://subs:subs@localhost:5435/subs?sslmode=disable

# Цели не связаны с файлами и их надо выполнять всегда
.PHONY: build run docker-build compose-up compose-down migrate-up

# Собираем локально бинарник
build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -o $(BIN_DIR)/$(APP) $(MAIN)

# Запускаем main.go
run:
	go run $(MAIN)

# Собираем Docker-образ
docker-build:
	docker build -t $(IMG) .

# Запускаем docker-compose.yml
compose-up:
	docker compose up --build -d

# Останавливаем и очищаем
compose-down:
	docker compose down -v

# Запускаем миграции
migrate-up:
	docker run --rm -v "$(PWD)/migrations:/migrations" --network host \
		migrate/migrate:4 -path=/migrations -database "$(DB_DSN)" up
