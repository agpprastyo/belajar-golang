include .env
export $(shell sed 's/=.*//' .env)



MIGRATE=go run github.com/golang-migrate/migrate/v4/cmd/migrate
#DB_URL=migrate -path ./migrations -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(127.0.0.1:3306)/$(MYSQL_DATABASE)?parseTime=true"
DB_URL=mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@127.0.0.1:3306/$(MYSQL_DATABASE)?sslmode=disable

# Targets
.PHONY: migrate-up migrate-down

migrate-up:
	migrate -database $(DB_URL) -path ./migrations up

migrate-down:
	$(MIGRATE) -database $(DB_URL) -path ./migrations down
