package internal

import (
	pb "auth-server/authService/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
	PostgresPool *pgxpool.Pool
}


func NewStruct() *AuthServer {
	return &AuthServer{}
}