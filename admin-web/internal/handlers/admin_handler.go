package handlers

import (
"net/http"
"admin-web/internal/response"
)


//Обработчик для страницы Admin:
func AdminHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		response.SendResponseError(w, http.StatusMethodNotAllowed, "Method not allowed!")
		return
	}


	// Проверяем сессию:
	// Для начала получаем токен и Cookie из запроса. При этом, если нету или ошибка, автоматически отправляем ошибку
	// token, err := auth.CheckSession(w, req)
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
	// 		return
	// 	}
	// 	response.SendResponseError(w, http.StatusInternalServerError, "Internal Server Error!")
	// 	return
	// }

	//Делаем запрос 
	
	response.SendResponseSuccess(w, http.StatusOK, "Welcome to admin-panel", nil)
	return
}
