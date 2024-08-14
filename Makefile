MIGRATION_FOLDER=$(CURDIR)/migrations
DOCKER_COMPOSE_FILE=docker-compose.yml
POSTGRES_SETUP_TEST ?= user=postgres password=postgres dbname=postgres host=postgres port=5432 sslmode=disable

.PHONY: docker-compose-up migration-up migration-down build docker-build

migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/${PROJECT_NAME} ./cmd/house-service

docker-build: build
	docker build -t house-service:latest .

docker-compose-up: docker-build
	docker-compose up