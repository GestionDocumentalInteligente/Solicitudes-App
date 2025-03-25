SHELL := /bin/bash

# Variables
ROOT_DIR := $(shell pwd)
VERSION := 1.0
BUILD_DIR := ${ROOT_DIR}/bin
CONFIG_DIR := ${ROOT_DIR}/config
SCRIPTS_DIR := ${ROOT_DIR}/scripts
EXAMPLES_DIR := ${ROOT_DIR}/examples
DOCKER_COMPOSE_FILE := ${CONFIG_DIR}/docker-compose.dev.yml

# Phony targets
.PHONY: all build run test docker-up docker-build docker-down docker-logs clean lint greeter-server-up greeter-server-build greeter-server-down greeter-server-logs greeter-client-up greeter-client-build greeter-client-down greeter-client-logs monitoring-up monitoring-build monitoring-down go-micro-up go-micro-build go-micro-down sg_auth-up sg_auth-build sg_auth-down sg_auth-logs

# Default target
all: build run

# Build target
build:
	@echo "Building the project..."
	@mkdir -p ${BUILD_DIR}
	# El APP_NAME debe definirse explícitamente
	APP_NAME=main go build -gcflags "all=-N -l" -o ${BUILD_DIR}/${APP_NAME} -ldflags "-X main.Version=${VERSION}" ${ROOT_DIR}/cmd/main.go

# Run target
run:
	@echo "Running the project..."
	@go run ${ROOT_DIR}/cmd/main.go

# Test target
test:
	@echo "Running tests..."
	@go test ./...

# Generic Docker commands
docker-up:
	@echo "Starting $(PROFILE) services in dev mode..."
	@chmod +x ${SCRIPTS_DIR}/entrypoint.sh
	APP_NAME=$(APP_NAME) docker compose -f ${DOCKER_COMPOSE_FILE} --profile $(PROFILE) up -d
	@$(MAKE) docker-logs PROFILE=$(PROFILE)

docker-build:
	@echo "Building $(PROFILE) services in dev mode..."
	@chmod +x ${SCRIPTS_DIR}/entrypoint.sh
	APP_NAME=$(APP_NAME) docker compose -f ${DOCKER_COMPOSE_FILE} --profile $(PROFILE) up --build -d
	@$(MAKE) docker-logs PROFILE=$(PROFILE)

docker-down:
	@echo "Stopping $(PROFILE) services in dev mode..."
	APP_NAME=$(APP_NAME) docker compose -f ${DOCKER_COMPOSE_FILE} --profile $(PROFILE) down --remove-orphans

docker-logs:
	@echo "Fetching logs for $(PROFILE) services..."
	APP_NAME=$(APP_NAME) docker compose -f ${DOCKER_COMPOSE_FILE} logs -f $(APP_NAME)

sg_auth-up:
	@$(MAKE) docker-up PROFILE=sg_auth APP_NAME=sg_auth

sg_auth-build:
	@$(MAKE) docker-build PROFILE=sg_auth APP_NAME=sg_auth

sg_auth-down:
	@$(MAKE) docker-down PROFILE=sg_auth APP_NAME=sg_auth

sg_auth-logs:
	@$(MAKE) docker-logs PROFILE=sg_auth APP_NAME=sg_auth

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -f ${BUILD_DIR}/${APP_NAME}

# Lint target
lint:
	@echo "Linting the project..."
	@golangci-lint run --config .golangci.yml --verbose
