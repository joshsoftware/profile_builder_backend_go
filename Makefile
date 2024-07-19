run: ## Run project on host machine
	go run cmd/main.go
	
clean: ## Clean database file for a fresh start
	rm test.db

test: ## Run all unit tests in the project
	go test -v ./...

test-cover: ## Run all unit tests in the project with test coverage
	go test -v ./... -covermode=count -coverprofile=coverage.out

html-cover: test-cover
	go tool cover -html="coverage.out"

migrate-up:
	go run internal/db/migrate.go up

migrate-down:
	go run internal/db/migrate.go down