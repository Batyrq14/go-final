.PHONY: help build up down logs clean dev-up dev-down migrate-up migrate-down test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker images
	docker-compose build

up: ## Start all services (production)
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View API logs
	docker-compose logs -f api

clean: ## Stop services and remove volumes
	docker-compose down -v

dev-up: ## Start development environment with hot reload
	docker-compose -f docker-compose.dev.yml up

dev-down: ## Stop development environment
	docker-compose -f docker-compose.dev.yml down

dev-logs: ## View development API logs
	docker-compose -f docker-compose.dev.yml logs -f api

migrate-up: ## Run database migrations
	docker-compose run --rm migrate -path /migrations -database "postgres://postgres:postgres@postgres:5432/qasynda?sslmode=disable" up

migrate-down: ## Rollback database migrations
	docker-compose run --rm migrate -path /migrations -database "postgres://postgres:postgres@postgres:5432/qasynda?sslmode=disable" down

test: ## Run tests
	go test ./...

test-coverage: ## Run tests with coverage
	go test -cover ./...

db-shell: ## Access PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d qasynda

