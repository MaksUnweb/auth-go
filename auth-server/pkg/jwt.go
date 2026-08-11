package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	JWT string 
	SessionID string
}

// Функция для генерации токена, исходя из собранных из сигнатуры claims
// Возвращает строку с токеном в формате header.claims.signature или ошибку
// Токен действует 1 час, соответственно и сессия длится 1 час
func GenerateToken(userID string, secret string) (*JWT, error) {

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"jti": sessionID, 
		"iat": time.Now().Unix(), 
		"exp": time.Now().Add(15 * time.Second).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	JWTStruct := &JWT{
		JWT: tokenString,
		SessionID: sessionID,
	}
	return JWTStruct, err
}


//Генерация UUID сессии:
func generateSessionID() (string, error){
	sessionID := make([]byte, 32)
	_, err := rand.Read(sessionID)
	if err != nil {
		log.Printf("Error generate sessionID: %v", err)
		return "", err
	}

	return hex.EncodeToString(sessionID), nil
}


func VerifyToken() {
	
}
