package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/message"
	authpb "abbassmortazavi/go-microservice/pkg/proto/auth"
	permissionpb "abbassmortazavi/go-microservice/pkg/proto/permission"
	rbacpb "abbassmortazavi/go-microservice/pkg/proto/rbac"
	rolepb "abbassmortazavi/go-microservice/pkg/proto/role"
	"abbassmortazavi/go-microservice/services/auth-service/grpc"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	"abbassmortazavi/go-microservice/services/auth-service/security"
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
	gcfg := global.Load()

	// ---- RabbitMQ ----
	message.Init()
	publisher := message.GetPublisher()

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
	tokenService := service2.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo, userRepo)
	authService := service2.NewAuthService(userRepo, hasher, tokenService, publisher)
	middlewares.Init(tokenService, authService)

	roleRepo := repository.NewRoleRepository(database.DB)
	permissionRepo := repository.NewPermissionRepository(database.DB)
	rbacRepo := repository.NewRBACRepository(database.DB)

	//authService := service2.NewAuthService(userRepo, hasher, tokenService, publisher)
	rbacService := service2.NewRBACService(userRepo, roleRepo, permissionRepo, rbacRepo)
	permissionService := service2.NewPermissionService(permissionRepo)
	roleService := service2.NewRoleService(roleRepo)

	authHandler := grpc.NewAuthHandler(authService)
	rbacHandler := grpc.NewRabcHandler(rbacService)
	permissionHandler := grpc.NewPermissionHandler(permissionService)
	roleHandler := grpc.NewRoleHandler(roleService)

	server := grpc2.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)
	rbacpb.RegisterRBACServiceServer(server, rbacHandler)
	permissionpb.RegisterPermissionServiceServer(server, permissionHandler)
	rolepb.RegisterRoleServiceServer(server, roleHandler)

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
