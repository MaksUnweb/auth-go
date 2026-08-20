package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "admin-web/authService/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Поддельная реализация клиента gRPC для тестирования:
type mockAuthClient struct {
	checkAuthFunc func(context.Context, *pb.CheckAuthRequest, ...grpc.CallOption) (*pb.CheckAuthResponse, error)
}

func(m *mockAuthClient) CheckAuth(ctx context.Context, req *pb.CheckAuthRequest, opts ...grpc.CallOption) (*pb.CheckAuthResponse, error) {
	return m.checkAuthFunc(ctx, req, opts...)
}



// Тестирования сценария, когда нету Cookie: 
func TestAuthMiddlewareNoCookie(t *testing.T) {
	client := &mockAuthClient{ 
		checkAuthFunc: func(ctx context.Context, car *pb.CheckAuthRequest, co ...grpc.CallOption) (*pb.CheckAuthResponse, error) {
			t.Fatal("CheckAuth must not be called!")
			return nil, nil
		},
	}


	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Next not be called! None cookie!")
	})

	handler := AuthMiddleware(next, client)

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)

	recorder := httptest.NewRecorder() 
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", recorder.Code)
	}
}

// Тестирования сценария, когда токен не валидный:
func TestAuthMiddlewareNoValideCookie(t *testing.T) {
	client := &mockAuthClient{
		checkAuthFunc: func(ctx context.Context, car *pb.CheckAuthRequest, co ...grpc.CallOption) (*pb.CheckAuthResponse, error) {
			return nil, status.Error(codes.Unauthenticated, "Unauthenticated!")
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Error! next func not be called!")
	})
	

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name: "session",
		Value: "invalid-token",
	})
	
	handler := AuthMiddleware(next, client)


	recorder := httptest.NewRecorder() 
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", recorder.Code)
	}
}


// Тестирование, когда произошла ошибка сервера: 
func TestAuthMiddlewareInternal(t *testing.T){
	client := &mockAuthClient{
		checkAuthFunc: func(ctx context.Context, car *pb.CheckAuthRequest, co ...grpc.CallOption) (*pb.CheckAuthResponse, error) {
			return nil, status.Error(codes.Internal, "Internal server error!")
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
		t.Fatal("Error! Next func not be called!")
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name: "session",
		Value: "some token...",
	})

	handler := AuthMiddleware(next, client)

	recorder := httptest.NewRecorder() 
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d", recorder.Code)
	}
}

// Тестирование, когда всё идиально: 
func TestAuthMiddlewareSuccess(t *testing.T) {
	client := &mockAuthClient{
		checkAuthFunc: func(ctx context.Context, car *pb.CheckAuthRequest, co ...grpc.CallOption) (*pb.CheckAuthResponse, error) {
			return &pb.CheckAuthResponse{
				IsAuth: true,
				Login: "Test_Login",
				AdminId: "adsd-13123-asfsa",
			}, nil
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name: "session",
		Value: "some token...",
	})

	handler := AuthMiddleware(next, client)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expect 200, got %d", recorder.Code)
	}
}
