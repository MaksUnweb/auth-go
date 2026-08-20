package handlers

import (
	pb "admin-web/authService/auth"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"admin-web/internal/models"
	"admin-web/internal/response"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type Loginer interface {
	Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginReply, error)
}


func LoginPost(w http.ResponseWriter, req *http.Request, client Loginer) {
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

	if input.Login == "" || input.Password == "" {
		response.SendResponseError(w, http.StatusBadRequest, "Bad request!")
		return
	}
	
	
	resp, err := client.Login(req.Context(), &pb.LoginRequest{Login: input.Login, Password: input.Password})
	if resp == nil && err == nil {
		response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
		return
	}
	
	if err != nil {
		code := status.Code(err)
	
		switch code{
			case codes.InvalidArgument:
				response.SendResponseError(w, http.StatusUnprocessableEntity, "Unprocessable Entity!")
				return
			case codes.Unauthenticated:
				response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
				return
			case codes.NotFound:
				response.SendResponseError(w, http.StatusNotFound, "Not Found!")
				return
			default:
				response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
				return 
		}
		
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
