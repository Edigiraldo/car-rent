.PHONY: build-db run-db build-server run-server

create-network:
	docker network create -d bridge cars-service-network

build-db:
	docker build . -f build/postgres/Dockerfile.local -t postgres_db-img

run-db:
	docker run -d --rm --network=cars-service-network -e POSTGRES_PASSWORD=postgres --name postgres_db postgres_db-img

run-debug-db:
	docker run -d --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 --name postgres_db postgres_db-img

build-server:
	docker build . -f build/server/Dockerfile.local -t car-rent-api-img

run-server:
	docker run -d --rm --network=cars-service-network -e ENVIRONMENT=local -p 5050:5050 --name car-rent-api car-rent-api-img

stop:
	docker stop car-rent-api postgres_db

compose-up:
	docker compose -f docker-compose.local.yml up -d

compose-down:
	docker compose -f docker-compose.local.yml down

compose-logs:
	docker compose -f docker-compose.local.yml logs -f

swagger-init:
	swag init --output doc/swagger -g cmd/api/main.go

run-tests:
	go test -cover ./...

all: create-network build-db run-db build-server run-server
