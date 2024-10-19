# Makefile

# Variables
MIGRATE=go run github.com/golang-migrate/migrate/v4/cmd/migrate
DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable

# Load environment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: migrate-up migrate-down

# Target to migrate up
migrate-up:
	migrate -database $(DB_URL) -path ./migrations up

# Target to migrate down
migrate-down:
	migrate -database $(DB_URL) -path ./migrations


# Target to run the application
run:
	go run ./cmd/api

# Target to build the application
build:
	go build -o ./bin/app ./cmd/api

# Target to run air for live reloading
air:
	air