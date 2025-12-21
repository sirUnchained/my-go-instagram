CONFIG_FILE ?= ./configs.json
MIGRATIONS_PATH := ./cmd/migrate
DB_ADDR := $(shell jq -r '.pg_db.addr' $(CONFIG_FILE))

.PHONY: run-dev
rub-dev:
	@ go run cmd/*.go

.PHONY: db-up
databases-up:
	@ docker compose up -d

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down