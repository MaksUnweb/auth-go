package internal

import (
	pb "auth-server/authService/auth"

	"auth-server/internal/auth"
	"auth-server/internal/repositories"
	"context"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (s AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) { 

	//Верифицируем токен и получаем его поля:
	jwtValues, err := auth.VerifyJWT(req.Token, s.JWTSecret)
	if err != nil {
		if err != auth.ErrInvalidToken {
			return nil, status.Error(codes.Internal, "Internal server error!")
		}
	}

	//Удаляем из Redis: 
	err = repositories.DelRedisSession(ctx, jwtValues.SessionID, s.RedisClient)
	if err != nil {
		if err != redis.Nil {
			return nil, status.Error(codes.Internal, "Internal server error!")
		}
	}
	
	return &pb.LogoutResponse{IsLogout: true}, nil
}
