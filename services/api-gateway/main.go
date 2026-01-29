package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/env"
	"abbassmortazavi/go-microservice/services/api-gateway/routes"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"abbassmortazavi/go-microservice/services/auth-service/repository"
	"abbassmortazavi/go-microservice/services/auth-service/service"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var (
	httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8085")
)

func main() {
	log.Println("Starting API Gateway")
	gcfg := global.Load()
	database.Connect()

	tokenRepo := repository.NewTokenRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	tokenService := service.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo, userRepo)
	middlewares.Init(tokenService)

	router := mux.NewRouter()
	routes.SetupRoutes(router)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Listening on ", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)
	case sig := <-shutdown:
		log.Println("Shutting down server...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down server: %v", err)
			server.Close()
		}
	}
}
