package auth

import (
	"log"
	"net/http"
)


//Функция для проверки сессии. 
// Получает данные из Cokie (JWT-токен)
func CheckSession(w http.ResponseWriter, req *http.Request) (string, error) {

	cookie, err := req.Cookie("session")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("Error parse cookie: %v", err)
		}
		log.Printf("Error: %v", err)
		return "", err
	}
	
	return cookie.Value, nil
	
}
