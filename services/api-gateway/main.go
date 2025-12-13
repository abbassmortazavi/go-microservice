package main

import (
	"abbassmortazavi/go-microservice/pkg/env"
	"log"
	"net/http"
)

var (
	httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8080")
)

func main() {
	log.Println("Starting API Gateway")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Preview called")
	})
}
