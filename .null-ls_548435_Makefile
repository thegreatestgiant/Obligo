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

test: 
	go test -v ./backend/...

reset-test: 
	docker compose rm -s -fv db-test
	docker compose up -d db-test

run:
	go run -C backend .

migrate-dev:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE is not set. Usage: make migrate-dev FILE=filename.sql"; \
		exit 1; \
	fi
	psql $(DB_DEV_URL) -f backend/schema/$(FILE)

db-wipe-dev:
	psql $(DB_DEV_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

db-setup-dev: db-wipe-dev
	psql $(DB_DEV_URL) -f backend/schema/000_official_look.sql

docker-build:
	docker buildx build --load -t obligo-app .

docker-run:
	docker run -p 8080:1234 \
		-e APP_PORT=1234 \
		-e APP_URL=http://localhost:8080 \
		obligo-app
