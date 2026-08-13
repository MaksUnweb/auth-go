package main

import (
	"log"
	
	"net/http"
	"admin-web/internal/handlers"
	"admin-web/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "admin-web/authService/auth"
)


func main() {

	//Создаю подключение к gRPC-серверу:
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

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
	

	// Основные маршруты. 
	// /login имеет 2 разных обработчика, так как помимо страндартного адреса с GET-параметром есть ещё и POST для обработки значений формы
	// http.HandleFunc("GET /{$}", handlers.HomeHandler)
	// http.HandleFunc("GET /login", handlers.LoginHandler)
	// http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
	// 		handlers.LoginPost(w, r, client)	
	// })
	// // http.HandleFunc("GET /admin", handlers.AdminHandler)
	//  http.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
	// 		 handlers.AdminHandler(w, r, &client)
	// 	})


	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Println("Запуск веб-сервера на порту :8080...")
	server.ListenAndServe()
}
