package handlers

import (
	"net/http"
	"admin-web/internal/response"
)

// Функция-обработчик для домашней страницы
func HomeHandler(w http.ResponseWriter , r *http.Request) {
	response.SendResponseSuccess(w, http.StatusOK, `Welcome to the test API for authentication. An example of authentication and authorization using a micro-server architecture is implemented here. 
		The services communicate using gRPC technology. 
		JWT is used as the token, and the session is additionally stored in Redis.`, nil)
}
