package main

// Основной сервер для аутентификации и авторизации.
// Обрабатывает основные запросы к базам данных, генерирует и верифицирует токены,
// работает с помощью gRPC, получает запросы с веб-сервера администратора, отправляет ответы
// При запуске сервера происходит миграция баз данных.
// Подход с использованием embed выбран так как проект является по большей степени примером,
// в полноценных production-решениях нужно использовать goose или golang-migrations

import (
	"context"
	"embed"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	pb "auth-server/authService/auth"
	"auth-server/internal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

//go:embed migrations/schema.sql
var schema embed.FS

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Внимание! Не получилось получить файл .env: %v", err)
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
		Addr:     rdbAddr,
		Password: rdbPass,
		DB:       intRdbDb,
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()


	// Запускаю миграции с текущим файлом:
	log.Println("Запускаю миграции...")
	data, err := schema.ReadFile("migrations/schema.sql")
	if err != nil {
		log.Fatalf("Ошибка чтения файла миграций: %v", err)
	}
	
	err = migrate(ctx, pool, string(data))
	if err != nil {
		log.Fatalf("Ошибка создания миграций: %v", err)
	}
	
	
	//Запускаем сервер:
	StartServer(ctx, pool, rdb, JWTSecret)
}


// Функция для создания миграции
// Миграции берутся из файла schema.sql с помощью embed
func migrate(ctx context.Context, pool *pgxpool.Pool, schema string) error {
	_, err := pool.Exec(ctx, schema)
	return err
}


func StartServer(ctx context.Context, pool *pgxpool.Pool, redis *redis.Client, JWTSecret string) {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя порта :50051: %v", err)
	}

	//Создаём экземпляр структуры сервера аутентификации:
	authServer := &internal.AuthServer{
		PostgresPool: pool,
		RedisClient:  redis,
		JWTSecret:    JWTSecret,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, authServer)

	go func() {
		log.Println("Запуск gRPC-сервера на порту :50051...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Ошибка запуска gRPC-сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Получил сигнал завершения, начинаю shutdown...")

	grpcServer.GracefulStop()
	log.Println("gRPC-сервер успешно остановлен!")

}
