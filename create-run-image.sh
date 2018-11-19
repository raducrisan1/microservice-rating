#run this locally so that you do not need to restore all the time the external dependencies (go get)
#of course, in a CI/CD environment, you need to change this approach 
docker rm $(docker ps -aqf "name=microservice-rating")
docker build -t local/microservice-rating .
docker tag local/microservice-rating gcr.io/itdays-201118/microservice-rating
docker run \
    --name microservice-rating \
    -e STOCKINFO_GRPC_ADDR='172.17.0.1:3001' \
    -e RATING_INTERVAL='10' \
    -e RABBITMQ_ADDR='amqp://guest:guest@172.17.0.1:5672/' \
    -p 3030:3030 \
    local/microservice-rating
