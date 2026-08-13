package main

// Основной сервер для аутентификации и авторизации.
// Обрабатывает основные запросы к базам данных, генерирует и верифицирует токены,
// работает с помощью gRPC, получает запросы с веб-сервера администратора, отправляет ответы

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"auth-server/internal"
	pb "auth-server/authService/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)



func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error load .env: %v", err)
	}

	postgresUrl := os.Getenv("DATABASE_URL")
	rdbAddr := os.Getenv("RDB_ADDR")
	rdbPass := os.Getenv("RDB_PASSWORD")
	rdbDb := os.Getenv("RDB_DB")
	JWTSecret := os.Getenv("JWT_SECRET")
	intRdbDb, err := strconv.Atoi(rdbDb)
	if err != nil {
		log.Fatalf("Errror parse RDB_DB environment: %v", err)
	}
	

	//Создаём пул подключения к базе данных Postgres:
	pool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к СУБД Postgres: %v", err)
	}

	//Создаём подключение к Redis:
	rdb := redis.NewClient(&redis.Options{
		Addr: rdbAddr,
		Password: rdbPass,
		DB: intRdbDb,
	})
	defer rdb.Close()

	//Проверяем работоспособность базы данных:
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error working with DB: %v", err)
	}

	redisStatus := rdb.Ping(context.Background())
	if redisStatus.Err() != nil {
		log.Printf("Error working with redis: %v", redisStatus.Err())
	}
	
	//Запускаем сервер:
	StartServer(pool, rdb, JWTSecret)

}


func StartServer(pool *pgxpool.Pool, redis *redis.Client, JWTSecret string) {

	lis, err := net.Listen("tcp", ":50051")	
	if err != nil {
		log.Fatalf("Ошибка создания слушателя порта :50051: %v", err)
	}

	//Создаём экземпляр структуры сервера аутентификации:
	authServer := &internal.AuthServer{
		PostgresPool: pool,
		RedisClient: redis,
		JWTSecret: JWTSecret,
	}


	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, authServer)

	log.Println("Запуск grpc-севрера...")
	grpcServer.Serve(lis)

}
