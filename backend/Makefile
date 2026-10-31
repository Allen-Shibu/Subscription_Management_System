.PHONY: db-up db-down migrate-up migrate-down run

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	./migration/migrate.sh up

migrate-down:
	./migration/migrate.sh down

run:
	go run ./cmd/api
