package main

import (
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"net/http"
)

func permissionGroup(mux *http.ServeMux) {
	authMiddleware := middlewares.GetMiddleware()
	mux.Handle("POST /create-permission", authMiddleware.AuthMiddleware(http.HandlerFunc(handelCreatePermission)))
	mux.Handle("POST /delete-permission", authMiddleware.AuthMiddleware(http.HandlerFunc(handelDeletePermission)))

}
