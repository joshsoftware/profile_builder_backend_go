# Makefile

.PHONY: run clean test test-cover html-cover migrate-up migrate-down new-migration

run: ## Run project on host machine
	go run cmd/main.go
	
clean: ## Clean database file for a fresh start
	rm -f test.db

test: ## Run all unit tests in the project
	go test -v ./...

test-cover: ## Run all unit tests in the project with test coverage
	go test -v ./... -covermode=count -coverprofile=coverage.out

html-cover: test-cover ## Generate HTML test coverage report
	go tool cover -html="coverage.out"

migrate-down: ## Roll back the last database migration
	@bash -c 'set -o allexport; source .env; migrate -database $$DB_MIGRATION -path internal/db/migrations down 1'

migrate-up: ## Apply database migrations
	@bash -c 'set -o allexport; source .env; migrate -database $$DB_MIGRATION -path internal/db/migrations up'

new-migration: ## Create a new migration
	@if [ -z "$(name)" ]; then \
		echo "Migration name not provided. Usage: make new-migration name=<new_migration_name>"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir internal/db/migrations -seq -format 20060102150405 $(name)
