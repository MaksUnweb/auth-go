package main

// Основной сервер для аутентификации и авторизации.
// Обрабатывает основные запросы к базам данных, генерирует и верифицирует токены,
// работает с помощью gRPC, получает запросы с веб-сервера администратора, отправляет ответы

import (
	"context"
	"log"
	"os"

	"auth-server/server"

	"github.com/jackc/pgx/v5/pgxpool"
)



func main() {

	postgresUrl := os.Getenv("DATABASE_URL")

	//Создаём пул подключения к базе данных Postgres:
	pool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к СУБД Postgres: %v", err)
	}

	//Проверяем работоспособность базы данных:
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error working with DB: %v", err)
	}

	//Запускаем сервер:
	server.StartServer(pool)

}



