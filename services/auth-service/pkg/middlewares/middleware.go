package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"
	"context"
	"log"
	"net/http"
	"strings"
)

type Middleware struct {
	authenticator service.TokenServiceInterface
}

func NewAuthMiddleware(authenticator service.TokenServiceInterface) *Middleware {
	return &Middleware{
		authenticator: authenticator,
	}
}

type contextKey string

const UserContextKey = contextKey("user")

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("No Authorization header found")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Invalid Authorization header")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		token, err := m.authenticator.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
