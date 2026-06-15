package main

import (
	"abbassmortazavi/go-microservice/services/notification-service/config"
	"abbassmortazavi/go-microservice/services/notification-service/worker"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("Starting notification-service")
	cfg := config.Load()
	emailWrrker := worker.NewEmailWorker(cfg)
	go emailWrrker.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
}
