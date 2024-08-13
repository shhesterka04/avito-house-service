MIGRATION_FOLDER=$(CURDIR)/migrations
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: docker-compose-up migration-up migration-down

docker-compose-up:
	docker-compose up

POSTGRES_SETUP_TEST ?= user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable


migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down


.PHONY: build
build: ## Build the Go binary
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/${PROJECT_NAME} ./cmd/house-service

.PHONY: docker-build
docker-build: build ## Build the Docker image
	@echo "Building Docker image..."
	docker build --build-arg APP_NAME=${PROJECT_NAME} -t ${IMAGE_REGISTRY}/${PROJECT_NAME}:${IMAGETAG} .