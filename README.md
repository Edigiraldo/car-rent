# Car Rent Microservice

docker build . -f build/server/Dockerfile -t car-rent-api
docker run -d --rm -p 5050:5050 --name cart-rent-api car-rent-api
docker stop car-rent-api
docker rmi car-rent-api
