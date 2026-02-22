package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AnyRoleMiddleware struct {
	authService service.AuthServiceInterface
}

func NewAnyRoleMiddleware(authService service.AuthServiceInterface) *AnyRoleMiddleware {
	return &AnyRoleMiddleware{
		authService: authService,
	}
}

func (rm *AnyRoleMiddleware) RequireAnyRole(requiredRoles []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := User(r.Context())
			log.Println("userID RequireAnyRole: ", strings.TrimSpace(user.ID))
			log.Println("user RequireAnyRole: ", user.Name)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			cleanID := strings.TrimSpace(user.ID)
			userID, err := strconv.ParseInt(cleanID, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID format", http.StatusBadRequest)
				return
			}
			res, err := rm.authService.CheckUseHasRole(r.Context(), userID, requiredRoles)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !res {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}
}
