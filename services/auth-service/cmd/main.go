package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/env"
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	"abbassmortazavi/go-microservice/services/auth-service/internal/domain/service"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/config"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/db"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/messaging"
	"abbassmortazavi/go-microservice/services/auth-service/internal/infrastructure/security"
	"abbassmortazavi/go-microservice/services/auth-service/internal/interface/grpc"
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
	rabbitmqURL := env.GetString("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	gcfg := global.Load()
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		cancel()
	}()

	lis, err := net.Listen("tcp", cfg.HTTP_ADDR)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	database.Connect()

	userRepo := db.NewUserRepository(database.DB)
	//db.Run(database.DB)

	hasher := security.NewBcryptHasher()
	tokenService := service.NewJWTSecret([]byte(gcfg.JWT_SECRET))

	// ---- RabbitMQ ----
	conn, ch := messaging.NewRabbitMQ(rabbitmqURL)
	defer conn.Close()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		cfg.USER_EXCHANGE,
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
	publisher := messaging.NewPublisher(ch, cfg.USER_EXCHANGE)

	authService := service.NewAuthService(userRepo, hasher, tokenService, publisher)

	authHandler := grpc.NewAuthHandler(authService)

	server := grpc2.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)

	log.Printf("starting Auth Service grpc server at %s", lis.Addr().String())

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Printf("shutting down Auth Service grpc server")
	server.GracefulStop()
}
