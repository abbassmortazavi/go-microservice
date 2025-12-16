package main

import (
	"abbassmortazavi/go-microservice/pkg/env"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8080")
)

func main() {
	log.Println("Starting API Gateway")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("this is the last change!!!")
		log.Println("thats ok")
	})

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
