APP_NAME=application-service

.PHONY: run build tidy docker-up docker-down docker-build docker-restart docker-logs docker-ps docker-clean

run:
	go run ./cmd

build:
	go build -o ./bin/$(APP_NAME) ./cmd

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose up -d --build

docker-restart:
	docker compose restart app

docker-logs:
	docker compose logs -f app

docker-ps:
	docker compose ps

docker-clean:
	docker compose down -v