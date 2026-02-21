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
	protected.HandleFunc("/assign-role-to-user", handlers.AssignRoleToUser).Methods("POST")
	protected.HandleFunc("/check-user-has-role", handlers.CheckUserHasRole).Methods("POST")
}
