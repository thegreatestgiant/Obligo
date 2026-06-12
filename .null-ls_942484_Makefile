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

run:
	go run -C backend .

migrate-dev: 
	psql $(DB_DEV_URL) -f backend/schema/$(FILE)

migrate-prod: 
	psql $(DB_PROD_URL) -f backend/schema/$(FILE)
