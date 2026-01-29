package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

func UserRoutes(mux *mux.Router) {
	authMiddleware := middlewares.GetMiddleware()
	protected := mux.PathPrefix("/api/v1/users").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)
	protected.HandleFunc("/me", handlers.GetUser).Methods("GET")

}
