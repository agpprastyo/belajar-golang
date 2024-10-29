SHELL := /bin/bash

include .env
export $(shell sed 's/=.*//' .env)

MIGRATE=go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate
DB_URL="mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(localhost:3306)/$(MYSQL_DATABASE)?multiStatements=true"

migrate-up:
	$(MIGRATE) -path ./migrations -database $(DB_URL) up

migrate-down:
	$(MIGRATE) -path ./migrations -database $(DB_URL) down

.PHONY: migrate-up migrate-down
