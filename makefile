.PHONY: start
start:
	@ go run cmd/*.go

.PHONY: databases
databases:
	@ docker compose up -d