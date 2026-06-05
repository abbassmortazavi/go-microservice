package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/service"
)

func Init(tokenService service.TokenServiceInterface, authService service.AuthServiceInterface) {
	globalMiddleware = NewAuthMiddleware(tokenService, authService)
}
