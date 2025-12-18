package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/config"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/db"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/security"
	"abbassmortazavi/go-microservice/services/auth-service/internal/interface/grpc"
	"log"
	"net"

	_ "github.com/lib/pq"
	grpc2 "google.golang.org/grpc"
)

func main() {
	gcfg := global.Load()
	cfg := config.Load()

	database.Connect()

	userRepo := db.NewUserRepository(database.DB)

	hasher := security.NewBcryptHasher()
	tokenService := service.NewJWTSecret([]byte(gcfg.JWT_SECRET))

	authService := service.NewAuthService(userRepo, hasher, tokenService)

	authHandler := grpc.NewAuthHandler(authService)

	server := grpc2.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)
	lis, err := net.Listen("tcp", cfg.HTTP_ADDR)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = server.Serve(lis)
	if err != nil {
		return
	}

}
