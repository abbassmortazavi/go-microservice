package main

import (
	"abbassmortazavi/go-microservice/pkg/utils"
	"abbassmortazavi/go-microservice/services/api-gateway/auth/requests"
	"abbassmortazavi/go-microservice/services/api-gateway/grpc_clients"
	"encoding/json"
	"log"
	"net/http"
)

func handelRegister(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()

	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := authService.Client.Register(ctx, req.ToProto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}
}

func handelLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("mmmmmmmmmmmmmmmmmmmmm")
	var req requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	ctx := r.Context()
	authService, err := grpc_clients.NewAuthServiceClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := authService.Client.Login(ctx, req.ToProto())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, res)
	if err != nil {
		return
	}
	log.Println("login success")
}
