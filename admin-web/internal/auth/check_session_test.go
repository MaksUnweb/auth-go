package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)



// Тест проверки и получения Cookie. Сценарий когда нету записи Cookie:
func TestGetSessionNoCookie(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)
	
	_, err := GetSession(req)

	if err != nil {
		if err != http.ErrNoCookie {
			t.Fatalf("Error! Expect ErrNoCookie, got %d", http.ErrNoCookie)
		}
	}
}

// Тест когда есть Cookie: 
func TestGetSessionSuccess(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin",
		nil,
	)

	req.AddCookie(&http.Cookie{ 
		Name: "session",
		Value: "some token...",
	})


	token, err := GetSession(req)
	if err != nil{
		t.Fatal("Error should not be called!")
	}
	
	if token != "some token..." {
		t.Fatalf("Expect token: \"some token...\", got %s", token)
	}
}
