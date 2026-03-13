.PHONY: up down run build logs ps

COMPOSE_FILE := deployments/docker-compose.yml

## Docker Compose — инфраструктура (TimescaleDB + MongoDB)

up:                ## Поднять контейнеры
	docker compose -f $(COMPOSE_FILE) up -d

down:              ## Остановить контейнеры
	docker compose -f $(COMPOSE_FILE) down

logs:              ## Логи контейнеров (follow)
	docker compose -f $(COMPOSE_FILE) logs -f

ps:                ## Статус контейнеров
	docker compose -f $(COMPOSE_FILE) ps

## Локальный запуск сервера

run:               ## Запустить сервер локально
	go run ./cmd/server

build:             ## Собрать бинарник
	go build -o bin/server ./cmd/server

help:              ## Показать справку
	@grep -E '^[a-z]+:.*##' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  %-14s %s\n", $$1, $$2}'
