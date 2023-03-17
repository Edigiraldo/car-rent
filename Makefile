.PHONY: build-db run-db build-server run-server

build-db:
	docker build . -f build/postgres/Dockerfile -t postgres_db-img

run-db:
	docker run -d --rm -e POSTGRES_PASSWORD=postgres --name postgres_db postgres_db-img

build-server:
	docker build . -f build/server/Dockerfile -t car-rent-api-img

run-server:
	docker run -d --rm -p 5050:5050 --name car-rent-api car-rent-api-img

all: build-db run-db build-server run-server