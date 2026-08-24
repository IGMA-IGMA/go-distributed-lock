.PHONY: run build test docker-up docker-down migrate

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

test-race:
	go test ./... -race

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

migrate:
	go run ./cmd/migrate
