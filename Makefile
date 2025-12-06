PROJECT_NAME := qasynda
PROTO_SRC := shared/proto
PROTO_OUT := shared/proto

.PHONY: proto clean docker-up docker-down

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $(PROTO_SRC)/*.proto

tidy:
	go mod tidy

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

reset-db:
	docker-compose down -v
	docker-compose up -d

run-gateway:
	go run ./services/gateway

run-user:
	go run ./services/user

run-marketplace:
	go run ./services/marketplace

run-chat:
	go run ./services/chat
