package handlers

import (
	pb "admin-web/authService/auth"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



// Поддельная реализация клиента для логина:
type mockLogin struct {
	LoginFunc func(context.Context, *pb.LoginRequest, ...grpc.CallOption) (*pb.LoginReply, error)
	LogoutFunc func(context.Context, *pb.LogoutRequest, ...grpc.CallOption) (*pb.LogoutResponse, error)
	CheckAuthFunc func(context.Context, *pb.CheckAuthRequest, ...grpc.CallOption) (*pb.CheckAuthResponse, error)
}


//  Функция, которая срабатывает при вызове теста. Эта функция возвращает LoginFunc, остальные функции возвращают nil, nil, так как не используются в этом тесте
func (m *mockLogin) Login(
    ctx context.Context,
    req *pb.LoginRequest,
    opts ...grpc.CallOption,
) (*pb.LoginReply, error) {
    return m.LoginFunc(ctx, req, opts...)
}


// Тест для сценария входа, когда метод POST вообще не отправил данных:
func TestLoginPostBadRequest(t *testing.T) {
	client := &mockLogin{
		LoginFunc: func(ctx context.Context, lr *pb.LoginRequest, co ...grpc.CallOption) (*pb.LoginReply, error) {
			t.Fatal("Should not be a Client call!")
			return nil, nil
		},
	}

	body := strings.NewReader(`{"": "", "": ""}`)
	req := httptest.NewRequest(
		http.MethodPost, 
		"/login", 
		body,
	)

	recorder := httptest.NewRecorder()

	LoginPost(recorder, req, client)
	
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", recorder.Code)
	}
}

// Тест для сценария, когда данные оказались не валидные:
func TestLoginPostBadRequestNotValide(t *testing.T) {
	client := &mockLogin{
		LoginFunc: func(ctx context.Context, lr *pb.LoginRequest, co ...grpc.CallOption) (*pb.LoginReply, error) {
			return nil, status.Error(codes.InvalidArgument, "Invalide Data!")
		},
	}

	body := strings.NewReader(`{"login": "vacsa!@#!$", "password": "@!#%$@#$"}`)
	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		body,
	)
	
	recorder := httptest.NewRecorder() 
	LoginPost(recorder, req, client)
if recorder.Code != http.StatusUnprocessableEntity{
		t.Fatalf("Expected 422, got %d", recorder.Code)
	}
}


// Тест для сценария, когда сервер вернул ошибку 500:
func TestLoginPostInternal(t *testing.T) {
	client := &mockLogin{
		LoginFunc: func(ctx context.Context, lr *pb.LoginRequest, co ...grpc.CallOption) (*pb.LoginReply, error) {
			return nil, status.Error(codes.Internal, "Internal server error!")
		},
	}

	body := strings.NewReader(`{"login": "Test", "password": "test_password"}`)
	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		body,
	)
	
	recorder := httptest.NewRecorder() 

	LoginPost(recorder, req, client)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d", recorder.Code)
	}
}


// Тест для сценария, когда 200 ок: 
func TestLoginPostSuccess(t *testing.T) {
	client := &mockLogin{
		LoginFunc: func(ctx context.Context, lr *pb.LoginRequest, co ...grpc.CallOption) (*pb.LoginReply, error) {
			return &pb.LoginReply{
				Token: "asdadasdasdaa",
			}, nil
		},
	}

	body := strings.NewReader(`{"login": "Test", "password": "test_password"}`)
	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		body,
	)
	
	recorder := httptest.NewRecorder() 

	LoginPost(recorder, req, client)
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", recorder.Code)
	}
}
