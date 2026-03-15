package middleware

import (
	"context"
	"net/http"
	"strings"

	"gophermart/internal/utils"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")

		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		userID, err := utils.ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int64, bool) {

	id, ok := ctx.Value(userIDKey).(int64)

	return id, ok
}
