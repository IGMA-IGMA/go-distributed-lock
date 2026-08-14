.PHONY: run build docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

migrate:
	go run ./cmd/migrate
