package handlers

import (
	"log"
	"net/http"
)

func CreatePermission(w http.ResponseWriter, r *http.Request) {
	log.Println("create permission")
}
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	log.Println("delete permission")
}
