package main

import (
	global "abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"
	"abbassmortazavi/go-microservice/pkg/env"
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
)

var (
	httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")
	gcfg := global.Load()
	database.Connect()

	tokenRepo := repository.NewTokenRepository(database.DB)
	tokenService := service.NewJwtAuthenticator(gcfg.JWT_SECRET, tokenRepo)
	middlewares.Init(tokenService)
	authMiddleware := middlewares.GetMiddleware()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("everything work perfectly!!!!!")
	})

	mux.Handle("POST /register", http.HandlerFunc(handelRegister))
	mux.Handle("POST /login", http.HandlerFunc(handelLogin))

	mux.Handle("GET /user/me", authMiddleware.AuthMiddleware(http.HandlerFunc(handelGetUser)))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
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
