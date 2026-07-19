include .env
export

up:
	docker compose up -d

down:
	docker compose down

downF:
	docker compose down -v

logs:
	docker compose logs -f

init-test-db:
	@docker compose exec db psql -U $(DB_USER) -d $(DEV_NAME) -tc \
		"SELECT 1 FROM pg_database WHERE datname = '$(TEST_NAME)'" | grep -q 1 \
		|| docker compose exec db psql -U $(DB_USER) -d $(DEV_NAME) -c \
		"CREATE DATABASE $(TEST_NAME);"

test: init-test-db
	go test -v ./backend/...

run:
	go run -C backend .

migrate-dev:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE is not set. Usage: make migrate-dev FILE=filename.sql"; \
		exit 1; \
	fi
	psql $(DB_URL) -f backend/schema/$(FILE)

db-wipe-dev:
	psql $(DB_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

db-setup-dev: db-wipe-dev
	psql $(DB_URL) -f backend/schema/000_official_look.sql

docker-build:
	docker buildx build --load -t obligo-app .

docker-run:
	docker compose -f created-compose.yml up

docker-down:
	docker compose -f created-compose.yml down
