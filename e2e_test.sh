#!/bin/bash
export LANG=ru_RU.utf8
export LC_ALL=ru_RU.utf8
echo "Запуск сервера..."
./start_server.sh & 
SERVER_PID=$!
sleep 5 

echo "Поиск комиксов по словам 'apple,doctor'..."
RESPONSE=$(curl -s "http://localhost:8080/pics?search="apple,doctor"")

if [[ "$RESPONSE" == *"apple a day"* ]]; then
echo "Тест пройден: найден комикс 'apple a day'"
else
echo "Тест не пройден: комикс 'apple a day' не найден"
exit 1
fi

echo "Остановка сервера..."
kill $SERVER_PID