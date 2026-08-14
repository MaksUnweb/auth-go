package handlers

import (
	"net/http"

	pb "admin-web/authService/auth"
	"admin-web/internal/auth"
	"admin-web/internal/response"

	"google.golang.org/grpc/codes"
)

// Обработчик для выхода из текущей сессии:
func Logout(w http.ResponseWriter, req *http.Request, client pb.AuthClient) {
	token, err := auth.GetSession(req)
	if err != nil{
		if err == http.ErrNoCookie {
			response.SendResponseSuccess(w, http.StatusNoContent, "No content!", nil)
			return
		}
		response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
		return
	}

	// Делаем запрос на gRPC-сервер аутентификации для выхода:
	_, err = client.Logout(req.Context(), &pb.LogoutRequest{Token: token})
	if err != nil {
		if err.Error() == codes.Unauthenticated.String() {
			response.SendResponseSuccess(w, http.StatusNoContent, "No content!", nil)
		}else{
			response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
		}

		DelCookie(w)
		return
	}

	// if !isLogout.IsLogout {
	// 	response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
	// 	return
	// }
	

	// Удаляем Cookie:
	DelCookie(w)
	
	response.SendResponseSuccess(w, http.StatusNoContent, "No content!", nil)
}


func DelCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Path: "/",
		MaxAge: -1,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
	})
}
