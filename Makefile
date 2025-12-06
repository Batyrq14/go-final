PROJECT_NAME := qasynda

.PHONY: clean docker-up docker-down

tidy:
	go mod tidy

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

reset-db:
	docker-compose down -v
	docker-compose up -d
	sleep 3
	make migrate-up

migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5433/qasynda?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5433/qasynda?sslmode=disable" down

run-gateway:
	DATABASE_URL="postgres://user:password@localhost:5433/qasynda?sslmode=disable" go run ./services/gateway

run-user:
	DATABASE_URL="postgres://user:password@localhost:5433/qasynda?sslmode=disable" go run ./services/user

run-marketplace:
	DATABASE_URL="postgres://user:password@localhost:5433/qasynda?sslmode=disable" go run ./services/marketplace

run-chat:
	DATABASE_URL="postgres://user:password@localhost:5433/qasynda?sslmode=disable" go run ./services/chat
