include .env
export

# --- Docker Commands --- #

# up: Start all services in detached mode and rebuild images
up:
	docker compose up -d --build

# watch: Start services and show logs 
watch:
	docker compose up --build

# down: Stop and remove containers
stop:
	docker compose down

# Rebuild everything from scratch without using cache
rebuild:
	docker compose build --no-cache
	docker compose up -d

# clean: Stop containers and remove volumes (Wipe everything)
clean:
	docker compose down -v

# --- Goose Migration Commands --- #

# goose-up: Run all pending migrations
goose-up:
	docker compose exec app goose up

# goose-down: Roll back a single migration
goose-down:
	docker compose exec app goose down

# goose-create: Create a new migration file (Usage: make goose-create name=add_users_table)
goose-create:
	docker compose exec app goose create $(name) sql
	sudo chown -R $(shell id -u):$(shell id -g) ./migrations

# goose-status: Check migration status
goose-status:
	docker compose exec app goose status

# goose-reset: Roll back all migrations
goose-reset:
	docker compose exec app goose reset

# db-flush: Drop everything and start migrations fresh (Dangerous!)
db-flush:
	docker compose exec app goose reset
	docker compose exec app goose up

.PHONY: up watch stop clean goose-up goose-down goose-create goose-status goose-reset db-flush
