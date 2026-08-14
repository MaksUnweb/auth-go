package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var ErrInvalidToken = errors.New("invalid token")


type JWT struct {
	JWT string 
	SessionID string
}


// Структура для хранения данных из JWT-токена
type JWTValues struct {
	AdminID string 
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
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	JWTStruct := &JWT{
		JWT: tokenString,
		SessionID: sessionID,
	}
	return JWTStruct, err
}


// Верификация токена:
func VerifyJWT(tokenString, secret string) (*JWTValues, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("Error Verivication token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}
	
	return &JWTValues{
		AdminID: claims.Subject,
		SessionID: claims.ID,
	}, nil
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
