package internal

import (
	pb "auth-server/authService/auth"
	"auth-server/internal/auth"
	"auth-server/internal/repositories"
	"context"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type AuthServer struct {
	pb.UnimplementedAuthServer
	PostgresPool *pgxpool.Pool
	RedisClient *redis.Client
	JWTSecret string
}


func NewStruct() *AuthServer {
	return &AuthServer{}
}



var dataRegex = regexp.MustCompile(`^[a-zA-Z0-9.@-_]+$`)

func (s AuthServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {

	//Валидация логина и пароля:
	if !validatePass(r.Password) || !validateLogin(r.Login) {
		return nil, status.Error(codes.InvalidArgument, "Login or password invalide!")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	//Получаю данные из бд об администраторе:
	admin,err := repositories.CheckInDb(r.Login, s.PostgresPool, ctx)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Error: %v", err)
			return nil, status.Error(codes.Internal, "Server error, please try again later!")
		}
		return nil, status.Error(codes.NotFound, "Incrorrect login or password!")
	}
	
	//Верифицируем пароль:
	err = auth.VerifyPass(r.Password, admin.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Incrorrect login or password!")
	}
	
	// Запускаем функцию для работы с сессиями:
	// Это действие генерирует JWT-токен и вносит в Redis данные JTI и userID как значение сессии: 
	token, err := auth.SessionPush(ctx, admin.AdminID, s.JWTSecret, s.RedisClient)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error!")
	}
	
	return &pb.LoginReply{Token: token}, nil
}


//Валидация пароля, максимум 100 символом, минимум 3 символа
func validatePass(pass string) bool {
	if len(pass) > 100 || len(pass) < 3{
		return false
	}

	return dataRegex.MatchString(pass)

}

//Валидация логина, максимум 150 символом, минимум 3 символа
func validateLogin(login string) bool {
	if len(login) > 150 || len(login) < 3{
		return false
	}

	return dataRegex.MatchString(login)
}
