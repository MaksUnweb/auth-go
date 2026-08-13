package handlers

import (
	pb "admin-web/authService/auth"
	"encoding/json"
	"net/http"
	"time"

	"admin-web/internal/models"
	"admin-web/internal/response"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func LoginPost(w http.ResponseWriter, req *http.Request, client pb.AuthClient) {
	if req.Method != "POST"	 {
		response.SendResponseError(w, http.StatusMethodNotAllowed, "Method not allowed!")
		return 
	}

	var input = &models.Login{}
	
	// Получаем POST-запрос и десериализуем данные в структуру из JSON: 
	if err := json.NewDecoder(req.Body).Decode(input); err != nil {
		response.SendResponseError(w, http.StatusBadRequest, "Bad request!")
		return
	}
	
	
	resp, err := client.Login(req.Context(), &pb.LoginRequest{Login: input.Login, Password: input.Password})

	if err != nil {
		st := status.Convert(err)
	
		if st.Code() == codes.InvalidArgument || st.Code() == codes.Unauthenticated || st.Code() == codes.NotFound {
			response.SendResponseError(w, http.StatusFound, "Login or password not allowed!")
			return
		}

		// http.Redirect(w, req, "/login?error="+"Ошибка работы сервера! Попробуйте позже!", http.StatusFound)
		response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
		return 
	}

	
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: resp.Token,
		Path: "/",
		Expires: time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
	})

	// http.Redirect(w, req, "/admin", http.StatusSeeOther)
	response.SendResponseSuccess(w, http.StatusOK, "Authorized success!", nil) 
	return
}
