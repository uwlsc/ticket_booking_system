include .env
export

MIGRATE=sql-migrate

create:
	@read -p  "What is the name of migration?" NAME; \
	${MIGRATE} new -config=config/dbconfig.yml $$NAME

up:
	docker compose --env-file .env -f deployments/docker/docker-compose.yml up

build:
	docker compose --env-file .env -f deployments/docker/docker-compose.yml up --build

down:
	docker compose --env-file .env -f deployments/docker/docker-compose.yml down -v

.PHONY: create up build down