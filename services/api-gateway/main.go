package main

import (
	"abbassmortazavi/go-microservice/pkg/env"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
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
	mux := http.NewServeMux()
	mux.HandleFunc("POST /test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("everything work perfectly!!!!!")
	})
	authMiddleware := middlewares.GetMiddleware()
	mux.Handle("POST /register", authMiddleware.AuthMiddleware(http.HandlerFunc(handelRegister)))
	mux.Handle("POST /login", authMiddleware.AuthMiddleware(http.HandlerFunc(handelLogin)))

	mux.HandleFunc("GET /user/:id", handelGetUser)

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
