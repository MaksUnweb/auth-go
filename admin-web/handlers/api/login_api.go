package api

import (
	pb "admin-web/authService/auth"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoginApi(w http.ResponseWriter, req *http.Request, client pb.AuthClient) {
	if req.Method != "POST"	 {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return 
	}

	err := req.ParseForm() 

	if err != nil {
		http.Error(w, "Error parse form", http.StatusBadRequest)
		return
	}

	login := req.Form.Get("login")
	password := req.Form.Get("password")
	
	response, err := client.Login(req.Context(), &pb.LoginRequest{Login: login, Password: password})

	if err != nil {
		st := status.Convert(err)
	
		if st.Code() == codes.InvalidArgument || st.Code() == codes.Unauthenticated || st.Code() == codes.NotFound {
			http.Redirect(w, req, "/login?error="+"Логин или пароль не верный!", http.StatusFound)
		}

		http.Redirect(w, req, "/login?error="+"Ошибка работы сервера! Попробуйте позже!", http.StatusFound)
		return 
	}

	log.Println(response)

	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}

