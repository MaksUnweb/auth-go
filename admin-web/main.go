package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"admin-web/internal/auth"
	"admin-web/internal/handlers"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "admin-web/authService/auth"
)

func main() {

	//Создаю подключение к gRPC-серверу:
	authServerAddr := os.Getenv("AUTH_SERVER_ADDR")
	if authServerAddr == "" {
		log.Printf("Ошибка получения адреса сервера аутентификации из переменных!")
		authServerAddr = "localhost:50051"
	}
	
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Ошибка создания gRPC-клиента: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)

	// Создаю мультиплексор для обработки маршрутов с middleware:
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginPost(w, r, client)
	})
	mux.Handle("GET /admin",
		auth.AuthMiddleware(http.HandlerFunc(handlers.AdminHandler), client),
	)
	mux.HandleFunc("POST /admin/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.Logout(w, r, client)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Запуск веб-сервера на порту :8080...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Получил сигнал завершения, начинаю shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при  shutdown: %v", err)
	}

	log.Println("Сервер остановлен!")
}
