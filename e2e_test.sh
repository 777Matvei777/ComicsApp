#!/bin/bash
container_id=$(docker run -d -e POSTGRES_PASSWORD=local -p 5432:5432 8e4fc9e18489)

sleep 5

./xkcd -c config.yaml &
server_pid=$!

sleep 5
jwt_token=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username": "Matvei", "password": "1234"}' http://localhost:8080/login)
if [ -z "$jwt_token" ]; then
    echo "Не удалось получить токен"
    docker stop $container_id
    docker rm $container_id
    exit 1
fi

curl -s -X POST -H "Authorization: $jwt_token" http://localhost:8080/update
sleep 2
curl -s http://localhost:8080/pics?search="apple,doctor"
# if echo "$response" | grep -q "apple a day"; then
#     echo "Тест пройден"
# else
#     echo "Тест не пройден"
# fi

kill %1
docker stop $container_id
docker rm $container_id