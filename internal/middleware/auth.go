package middleware

import (
	"net/http"
	"strings"
)

var validTokens = map[string]string{
	"client_token":    "client",
	"moderator_token": "moderator",
}

func AuthMiddleware(requiredType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userType, valid := validTokens[token]
			if !valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if requiredType == "moderator" && userType != "moderator" {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
