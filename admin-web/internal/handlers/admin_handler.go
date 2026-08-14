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

	
	response.SendResponseSuccess(w, http.StatusOK, "Welcome to admin-panel", nil)
	return
}
