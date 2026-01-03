package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/env"
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/config"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/messaging"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/security"
	"abbassmortazavi/go-microservice/services/auth-service/internal/interface/grpc"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	service2 "abbassmortazavi/go-microservice/services/auth-service/service"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	grpc2 "google.golang.org/grpc"
)

func main() {
	log.Println("Starting service Auth Service...")
	rabbitmqURL := env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	gcfg := global.Load()
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	database.Connect()

	userRepo := repository.NewUserRepository(database.DB)
	tokenRepo := repository.NewTokenRepository(database.DB)

	hasher := security.NewBcryptHasher()
	//tokenService := service.NewJWTSecret([]byte(gcfg.JWT_SECRET))
	tokenService := service2.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo)
	middlewares.Init(tokenService)

	// ---- RabbitMQ ----
	conn, ch := messaging.NewRabbitMQ(rabbitmqURL)
	defer conn.Close()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		cfg.UserExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	publisher := messaging.NewPublisher(ch, cfg.UserExchange)

	authService := service2.NewAuthService(userRepo, hasher, tokenService, publisher)

	authHandler := grpc.NewAuthHandler(authService)

	server := grpc2.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)

	// Run server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("starting Auth Service grpc server at %s", lis.Addr().String())
		if err := server.Serve(lis); err != nil {
			serverErr <- err
		}
	}()

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErr:
		log.Printf("gRPC server error: %v", err)
		cancel()
	case sig := <-sigchan:
		log.Printf("Received signal: %v", sig)
		cancel()
	case <-ctx.Done():
		// Context canceled by other means
	}

	log.Printf("shutting down Auth Service grpc server")
	server.GracefulStop()
}
