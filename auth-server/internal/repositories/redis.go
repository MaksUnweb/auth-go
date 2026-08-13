package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)


// Функция для проверки наличия сессии в Redis:
func CheckRedisSession(ctx context.Context, sessionID, userID string, rdb *redis.Client) error {
	val, err := rdb.Get(ctx, "session:"+sessionID).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("Session not exists!")
		}
		log.Printf("Error getting values from redis: %v", err)
	}

	if val != userID {
		return fmt.Errorf("Session not exists!")
	}
	
	return nil
}
