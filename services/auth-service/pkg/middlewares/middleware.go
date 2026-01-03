package middlewares

import (
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/service"
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

func CreateGRPCMiddleware(authClient authpb.AuthServiceClient) *GRPCAuthMiddleware {
	return &GRPCAuthMiddleware{
		authClient: authClient,
	}
}

type GRPCAuthMiddleware struct {
	authClient authpb.AuthServiceClient
}

type contextKey string

const UserContextKey = contextKey("user")

/*func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
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
		ctx := context.Background()
		claims, err := m.authenticator.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// Add user info to context
		ctx = context.WithValue(ctx, UserContextKey, map[string]interface{}{
			"user_id": claims.UserID,
			"name":    claims.Name,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}*/

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
