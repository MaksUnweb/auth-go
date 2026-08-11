package internal

import (
	pb "auth-server/authService/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
