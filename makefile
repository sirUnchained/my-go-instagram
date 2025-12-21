.PHONY: dev
dev:
	@ go run cmd/*.go

.PHONY: databases
databases:
	@ docker compose up -d