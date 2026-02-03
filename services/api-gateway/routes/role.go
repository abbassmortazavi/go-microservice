package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

func RoleRoutes(mux *mux.Router) {
	authMiddleware := middlewares.GetMiddleware()
	protected := mux.PathPrefix("/api/v1/roles").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)

	protected.HandleFunc("/create", handlers.Create).Methods("POST")
	protected.HandleFunc("/{id}", handlers.Delete).Methods("DELETE")
	protected.HandleFunc("/lists", handlers.List).Methods("GET")
	protected.HandleFunc("/{id}", handlers.Get).Methods("GET")
}
