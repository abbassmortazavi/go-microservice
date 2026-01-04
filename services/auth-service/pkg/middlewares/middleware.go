package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/interfaces/service_interface"
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"log"
	"net/http"
	"strings"
)

type Middleware struct {
	authenticator service_interface.TokenServiceInterface
}

func NewAuthMiddleware(authenticator service_interface.TokenServiceInterface) *Middleware {
	return &Middleware{
		authenticator: authenticator,
	}
}

type contextKey string

const UserContextKey = contextKey("user")

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if m == nil {
			log.Println("ERROR: Middleware is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if m.authenticator == nil {
			log.Println("ERROR: Token authenticator is nil")
			http.Error(w, "Authentication service unavailable", http.StatusServiceUnavailable)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("No Authorization header found")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
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

func (m *Middleware) GetUserFromContext(ctx context.Context) *service.Claims {
	user, _ := ctx.Value(UserContextKey).(*service.Claims)
	return user
}
