package server

import (
	"log"
	"net"

	"auth-server/internal"
	pb "auth-server/authService/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)



func StartServer(pool *pgxpool.Pool) {

	lis, err := net.Listen("tcp", ":50051")	
	if err != nil {
		log.Fatalf("Ошибка создания слушателя порта :50051: %v", err)
	}

	//Создаём экземпляр структуры сервера аутентификации:
	authServer := &internal.AuthServer{
		PostgresPool: pool,
	}


	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, authServer)

	log.Println("Запуск grpc-севрера...")
	grpcServer.Serve(lis)

}



// func (s *AuthServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {

// 	//Валидация логина и пароля:
// 	if !validatePass(r.Password) || !validateLogin(r.Login) {
// 		return nil, status.Error(codes.InvalidArgument, "Логин или пароль не валидный!")
// 	}

// 	return &pb.LoginReply{Token: "test-token!!wqe"}, nil
// }


// //Валидация пароля, максимум 100 символом, минимум 3 символа
// func validatePass(pass string) bool {
// 	if len(pass) > 100 || len(pass) < 3{
// 		return false
// 	}

// 	return dataRegex.MatchString(pass)

// }

// //Валидация логина, максимум 150 символом, минимум 3 символа
// func validateLogin(login string) bool {
// 	if len(login) > 150 || len(login) < 3{
// 		return false
// 	}

// 	return dataRegex.MatchString(login)
// }


// //Функция проверки в базе данных:
// func checkInDb(login string, pool *pgxpool.Pool)  {

// }
