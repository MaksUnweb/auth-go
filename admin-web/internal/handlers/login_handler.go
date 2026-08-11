package handlers

import (
	"net/http"
	"admin-web/internal/response"
)


func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		response.SendResponseError(w, http.StatusMethodNotAllowed, "Method not allowed!")
		return
	}
	
	response.SendResponseSuccess(w, http.StatusOK, "Login page! To do login into admin-panel, make a POST-request to /login with your login and password!", nil)
	return
}
