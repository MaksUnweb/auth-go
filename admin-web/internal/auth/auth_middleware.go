package auth

import (
	"log"
	"net/http"

	pb "admin-web/authService/auth"
	"admin-web/internal/response"

	"google.golang.org/grpc/codes"
)

func AuthMiddleware(next http.Handler, client pb.AuthClient) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		cookie, err := req.Cookie("session")
		if err != nil {
			if err != http.ErrNoCookie {
				log.Printf("Error parse cookie: %v", err)
			}
			
			response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
			return
		}

		// Делаем запрос на сервер аутентификации для проверки токена:
		resp, err := client.CheckAuth(req.Context(), &pb.CheckAuthRequest{Token: cookie.Value})
		if err != nil{
			log.Printf("Error CheckAuth: %v", err)
			if err.Error() == codes.Internal.String(){
				response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
				return
			}
			response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
			return
		}
		
		if !resp.IsAuth  {
			response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
			return
		}
		
		next.ServeHTTP(w, req)
	})
}
