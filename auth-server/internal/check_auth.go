package internal

import (
	pb "auth-server/authService/auth"
	"auth-server/internal/repositories"
	"context"

	"auth-server/internal/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (s AuthServer) CheckAuth(ctx context.Context, req *pb.CheckAuthRequest) (*pb.CheckAuthResponse, error) { 

	//Верифицируем токен:
	jwtValues, err := auth.VerifyJWT(req.Token, s.JWTSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error()) 
	}
	
	err = repositories.CheckRedisSession(ctx, jwtValues.SessionID, jwtValues.AdminID, s.RedisClient)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}


	//Проверяем есть ли такой пользователь в бд:
	login, err := repositories.CheckByAdminID(ctx, jwtValues.AdminID, s.PostgresPool)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	
	
	//Отправляем ответ. Здесь AdminID и Login администратора можно использовать для каких-либо целей в админ панели:
	return &pb.CheckAuthResponse{
		IsAuth: true,
		AdminId: jwtValues.AdminID,
		Login: login,
	}, nil
}
