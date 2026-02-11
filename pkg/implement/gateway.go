package implement

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/message"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	"abbassmortazavi/go-microservice/services/auth-service/security"
	"abbassmortazavi/go-microservice/services/auth-service/service"
)

func Implement() {
	gcfg := global.Load()
	database.Connect()
	hasher := security.NewBcryptHasher()
	// ---- RabbitMQ ----
	message.Init()
	publisher := message.GetPublisher()

	tokenRepo := repository.NewTokenRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	tokenService := service.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo, userRepo)
	authService := service.NewAuthService(userRepo, hasher, tokenService, publisher)
	middlewares.Init(tokenService, authService)
}
