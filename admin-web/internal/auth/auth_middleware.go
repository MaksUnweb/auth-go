package auth

import (
	"context"
	"log"
	"net/http"

	pb "admin-web/authService/auth"
	"admin-web/internal/response"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckAuthInt interface {
	CheckAuth(ctx context.Context, in *pb.CheckAuthRequest, opts ...grpc.CallOption) (*pb.CheckAuthResponse, error)
}


func AuthMiddleware(next http.Handler, client CheckAuthInt) http.Handler{
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
			if status.Code(err) == codes.Unauthenticated{
				response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
				return
			}
			response.SendResponseError(w, http.StatusInternalServerError, "Internal server error!")
			return
		}
		
		if !resp.IsAuth  {
			response.SendResponseError(w, http.StatusUnauthorized, "Unauthorized!")
			return
		}
		
		next.ServeHTTP(w, req)
	})
}
