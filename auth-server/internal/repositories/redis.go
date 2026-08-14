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


// Функция для удаления сессии из Redis: 
func DelRedisSession(ctx context.Context, sessionID string, rdb *redis.Client) error {
	key := "session:" + sessionID
	err := rdb.Del(ctx, key).Err() 
	if err != nil {
		if err != redis.Nil {
			log.Printf("Error request from Redis: %v", err)
		}
	}

	return err
}
