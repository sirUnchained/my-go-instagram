.PHONY: lets-gooooo
start:
	@ go run cmd/*.go

.PHONY: databases-be-ready
databases:
	@ docker compose up -d