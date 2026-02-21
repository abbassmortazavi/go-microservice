package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

func PermissionRoutes(mux *mux.Router) {
	authMiddleware := middlewares.GetMiddleware()
	anyRoleMiddleware := middlewares.GetAnyRoleMiddleware()
	protected := mux.PathPrefix("/api/v1/permissions").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)

	protected.HandleFunc("/create", handlers.CreatePermission).Methods("POST")
	protected.HandleFunc("/{id}", anyRoleMiddleware.RequireAnyRole([]string{"admin", "superadmin"})(handlers.DeletePermission)).Methods("DELETE")
	protected.HandleFunc("/lists", handlers.ListPermissions).Methods("GET")
	protected.HandleFunc("/{id}", handlers.GetPermission).Methods("GET")

}
