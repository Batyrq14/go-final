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

run-gateway:
	go run services/gateway/main.go

run-user:
	go run services/user/main.go

run-marketplace:
	go run services/marketplace/main.go

run-chat:
	go run services/chat/main.go
