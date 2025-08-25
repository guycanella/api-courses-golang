PKG := ./cmd/api
BIN := bin/api
SWAG := $(shell go env GOPATH)/bin/swag

.PHONY: run build tidy \
        test-handlers test-one test-race test-cover show-test-coverage \
        ps up down \
        migrate migrate-test \
        seed seed-test \
        swag

run:
	go run $(PKG)

build:
	mkdir -p bin
	go build -o $(BIN) $(PKG)

tidy:
	go mod tidy

ps:
	docker-compose ps

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	go run ./cmd/migrate $(ARGS)

migrate-test:
	$(MAKE) migrate ARGS="--test"

seed:
	go run ./cmd/seed $(ARGS)

seed-test:
	$(MAKE) seed ARGS="--test"

swag:
	$(SWAG) init -g cmd/api/main.go -o internal/docs

test-handlers:
	go test -v ./internal/handlers -count=1

test-race:
	go test -race ./internal/handlers

test-cover:
	go test -cover -coverprofile=coverage.out ./internal/handlers
	go tool cover -func=coverage.out

test-one:
	go test -v ./internal/handlers -run $(NAME) -count=1

show-test-coverage:
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html