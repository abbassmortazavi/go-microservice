package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"abbassmortazavi/go-microservice/services/auth-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

func PermissionRoutes(mux *mux.Router) {
	authMiddleware := middlewares.GetMiddleware()
	protected := mux.PathPrefix("/api/v1/permissions").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)

	protected.HandleFunc("POST /", handlers.CreatePermission)
	protected.HandleFunc("DELETE /{id}", handlers.DeletePermission)
	/*permGroup.Handle("GET /", handlers.ListPermissions)
	permGroup.Handle("GET /{id}", handlers.GetPermission)*/
}
