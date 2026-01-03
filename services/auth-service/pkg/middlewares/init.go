package middlewares

import "abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"

func Init(tokenService service.TokenServiceInterface) {
	globalMiddleware = NewAuthMiddleware(tokenService)
}
