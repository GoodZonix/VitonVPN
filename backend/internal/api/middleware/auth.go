package middleware

import (
	"context"
	"net/http"
	"strings"

	"vpn-startup/backend/internal/auth"
)

func JWT(jwt *auth.JWT) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(header, prefix) {
				http.Error(w, `{"error":"invalid authorization"}`, http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(header, prefix)
			claims, err := jwt.Verify(token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
