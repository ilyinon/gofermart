package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gophermart/internal/utils"
)

func TestAuthMiddleware(t *testing.T) {

	jwt := utils.NewJWTManager("test-secret")
	token, _ := jwt.GenerateToken(1)

	middleware := NewAuthMiddleware(jwt)
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, ok := GetUserID(r.Context())

		if !ok || id != 1 {
			t.Fatal("user id missing")
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
}
