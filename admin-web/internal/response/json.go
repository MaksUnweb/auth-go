package response

import (
	"encoding/json"
	"net/http"
)

//Функция для отправки ответов 200
//Принимает:
// w - заголовки запроса
// status - статус код
// message - сообщение 
// data - любые данные (могут быть nil)
func SendResponseSuccess(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	 json.NewEncoder(w).Encode(map[string]any{
		"message": message, 
		"data": data,
	})
}


// Функция для отправки ответов-ошибки:
// Принимает:
// w - заголовки запроса
// status - статус код
// message - сообщение 
func SendResponseError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	 json.NewEncoder(w).Encode(map[string]any{
		"message": message,
		})
}
