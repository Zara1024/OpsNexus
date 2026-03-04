APP_NAME=opsnexus-server

.PHONY: build run test lint tidy migrate-up migrate-down

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path ../migrations -database "postgres://postgres:postgres@127.0.0.1:5432/opsnexus?sslmode=disable" up

migrate-down:
	migrate -path ../migrations -database "postgres://postgres:postgres@127.0.0.1:5432/opsnexus?sslmode=disable" down
