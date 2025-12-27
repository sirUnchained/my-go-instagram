CONFIG_FILE ?= ./configs.json
MIGRATIONS_PATH := ./cmd/migrate
DB_ADDR := $(shell jq -r '.pg_db.addr' $(CONFIG_FILE))
SERVER_PATH := ../../cmd/server/

.PHONY: run-dev
run-dev:
	@ go run cmd/server/*.go

.PHONY: db-up
databases-up:
	@ docker compose up -d

.PHONY: db-down
databases-down:
	@ docker compose down

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down

.PHONY: generate-docs
generate-docs:
	@swag init -g ./cmd/server/main.go