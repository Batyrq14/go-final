#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Checking prerequisites...${NC}"

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go.${NC}"
    exit 1
fi

# Add GOPATH/bin to PATH
export PATH=$(go env GOPATH)/bin:$PATH

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install Docker.${NC}"
    exit 1
fi

# Check Protoc
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Protobuf Compiler (protoc) is not installed.${NC}"
    echo "Please install it using Homebrew:"
    echo -e "${GREEN}brew install protobuf${NC}"
    exit 1
fi

# Check Protoc Plugins
if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${GREEN}Installing Go Protobuf plugins...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

echo -e "${GREEN}Generating Proto definitions...${NC}"
make proto

echo -e "${GREEN}Tidying dependencies...${NC}"
go mod tidy

echo -e "${GREEN}Starting Infrastructure (Docker)...${NC}"
make docker-up

echo -e "${GREEN}Setup Complete!${NC}"
echo "You can now run the services in separate terminals using:"
echo "  make run-user"
echo "  make run-marketplace"
echo "  make run-chat"
echo "  make run-gateway"
