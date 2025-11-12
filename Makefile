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




