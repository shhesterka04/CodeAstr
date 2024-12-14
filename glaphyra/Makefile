MIGRATION_FOLDER=$(CURDIR)/migrations
DOCKER_COMPOSE_FILE=docker-compose.yml
POSTGRES_SETUP_TEST ?= user=test password=test dbname=test host=database port=5432 sslmode=disable

.PHONY: docker-compose-up migration-init migration-up migration-down build docker-build run

migration-up:
	docker-compose exec service goose -dir /app/migrations postgres "$(POSTGRES_SETUP_TEST)" up


migration-down:
	docker-compose exec service goose -dir /app/migrations postgres "$(POSTGRES_SETUP_TEST)" down


build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/glaphyra ./cmd/glaphyra

docker-build: build
	docker build -t glaphyra:latest .

docker-compose-up: docker-build
	docker-compose up -d

run: docker-compose-up migration-up
	docker-compose logs -f