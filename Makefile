SHELL := /bin/sh

.PHONY: build run test up down logs docker-build docs

build:
	go build -o bin/library-http ./cmd/library-http

run:
	go run ./cmd/library-http

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f app

docker-build:
	docker build -t library-http:local .

# create new db migration file
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migration-create NAME=add_books_table"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	FILENAME="migrations/$${TIMESTAMP}_$(NAME).sql"; \
	echo "-- +goose Up" > $$FILENAME; \
	echo "" >> $$FILENAME; \
	echo "-- +goose Down" >> $$FILENAME; \
	echo "Created new migration file: $$FILENAME"





