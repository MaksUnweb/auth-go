package auth

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Функция для создания сессий
// Запускает: генерацию токена, постановку сессии в Redis:
func SessionPush(ctx context.Context, userID string, secret string, rdb *redis.Client) (string, error) {
	JWTStruct, err := GenerateToken(userID, secret)
	if err != nil {
		return "", err
	}
	//Разбиваем полученную структуру на токен и sessionID: 
	token, sessionID := JWTStruct.JWT, JWTStruct.SessionID
	
	//Вставляем сессию в Redis:
	err = InsertSession(sessionID, userID, rdb, ctx)
	if err != nil {
		return "", err
	}

	return token, nil
}


// Функция для внесения сессии в Redis:
// Принимает JTI, который и является полезной нагрузкой сессии
func InsertSession(JTI, userID string, rdb *redis.Client, ctx context.Context) error {
	status := rdb.Set(ctx, "session:"+ JTI, userID, 1*time.Hour)
	if status.Err() != nil {
	 	log.Printf("Error set value into Redis: %v", status.Args()...)
		return status.Err()
	}
	return nil
}
