package handlers

import (
	"net/http"
)

// Функция-обработчик для домашней страницы
func HomeHandler(w http.ResponseWriter , r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>Hello, Maks</h1>"))
}
