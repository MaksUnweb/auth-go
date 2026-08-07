package main

import (
	"log"
	"net/http"

	"admin-web/handlers"
	"admin-web/handlers/api"

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

	

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/login-api", func(w http.ResponseWriter, r *http.Request) {
			api.LoginApi(w, r, client)	
	})
	http.HandleFunc("/admin", handlers.AdminHandler)

	log.Println("Запуск веб-сервера...")
	http.ListenAndServe(":8080", nil)
}

