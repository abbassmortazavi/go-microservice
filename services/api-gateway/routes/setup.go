package routes

import (
	"abbassmortazavi/go-microservice/services/api-gateway/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {

	public := router.PathPrefix("/api/v1").Subrouter()
	public.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		log.Println("everything work perfectly!!!!!")
	})
	public.HandleFunc("/register", handlers.Register).Methods("POST")
	public.HandleFunc("/login", handlers.Login).Methods("POST")

	// Register all service routes
	PermissionRoutes(router)
	UserRoutes(router)
	RoleRoutes(router)
}
