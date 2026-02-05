package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

func RBACRoutes(mux *mux.Router) {
	authMiddleware := middlewares.GetMiddleware()
	protected := mux.PathPrefix("/api/v1/rbac").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)

	protected.HandleFunc("/assign-permission-to-role", handlers.AssignPermissionToRole).Methods("POST")
	protected.HandleFunc("/{id}", handlers.Delete).Methods("DELETE")
	protected.HandleFunc("/lists", handlers.List).Methods("GET")
	protected.HandleFunc("/{id}", handlers.Get).Methods("GET")
}
