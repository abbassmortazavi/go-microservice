package middlewares

import (
	"abbassmortazavi/go-microservice/services/auth-service/service"
)

func Init(tokenService service.TokenServiceInterface) {
	globalMiddleware = NewAuthMiddleware(tokenService)
}
