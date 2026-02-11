package main

import (
	"abbassmortazavi/go-microservice/services/auth-service/pkg/implement"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting service Auth Service...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	implement.Init()
	// Setup signal handling
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := implement.GetServer()

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
