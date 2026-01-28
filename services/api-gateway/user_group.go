package main

import (
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"
	"net/http"
)

func userGroup(mux *http.ServeMux) {
	authMiddleware := middlewares.GetMiddleware()

	mux.Handle("GET /user/me", authMiddleware.AuthMiddleware(http.HandlerFunc(handelGetUser)))

}
