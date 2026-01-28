package main

import (
	"log"
	"net/http"
)

func publicGroup(mux *http.ServeMux) {
	mux.HandleFunc("POST /test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("everything work perfectly!!!!!")
	})
	mux.Handle("POST /register", http.HandlerFunc(handelRegister))
	mux.Handle("POST /login", http.HandlerFunc(handelLogin))
}
