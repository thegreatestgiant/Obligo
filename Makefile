include .env
export

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@grep -h -E '^[a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Start docker compose services in the background
	docker compose up -d

down: ## Stop docker compose services
	docker compose down

downF: ## Stop docker compose services and remove volumes
	docker compose down -v

logs: ## Tail docker compose logs
	docker compose logs -f

init-test-db: ## Initialize the test database if it doesn't exist
	@docker compose exec db psql -U $(DB_USER) -d $(DEV_NAME) -tc \
		"SELECT 1 FROM pg_database WHERE datname = '$(TEST_NAME)'" | grep -q 1 \
		|| docker compose exec db psql -U $(DB_USER) -d $(DEV_NAME) -c \
		"CREATE DATABASE $(TEST_NAME);"

test: init-test-db ## Run backend tests
	go test -v ./backend/...

run: ## Run the backend server locally
	go run -C backend .

migrate-dev: ## Run a specific migration file (e.g. make migrate-dev FILE=filename.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE is not set. Usage: make migrate-dev FILE=filename.sql"; \
		exit 1; \
	fi
	psql $(DB_URL) -f backend/schema/$(FILE)

db-wipe-dev: ## Wipe the dev database schema
	psql $(DB_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

db-setup-dev: db-wipe-dev ## Setup the dev database schema from scratch
	psql $(DB_URL) -f backend/schema/000_official_look.sql

docker-build: ## Build the production docker image
	docker buildx build --load -t obligo-app .

docker-run: ## Run the production docker compose setup
	docker compose -f created-compose.yml up

docker-down: ## Stop the production docker compose setup
	docker compose -f created-compose.yml down
