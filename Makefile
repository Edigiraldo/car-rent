.PHONY: build-db run-db build-server run-server

create-network:
	docker network create -d bridge cars-service-network

build-db:
	docker build . -f build/postgres/Dockerfile -t postgres_db-img

run-db:
	docker run -d --rm --network=cars-service-network -e POSTGRES_PASSWORD=postgres --name postgres_db postgres_db-img

build-server:
	docker build . -f build/server/Dockerfile -t car-rent-api-img

run-server:
	docker run -d --rm --network=cars-service-network -p 5050:5050 --name car-rent-api car-rent-api-img

stop:
	docker stop car-rent-api postgres_db

all: create-network build-db run-db build-server run-server