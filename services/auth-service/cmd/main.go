package main

import (
	"abbassmortazavi/go-microservice/pkg/config"
	"abbassmortazavi/go-microservice/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	database.Connect()
}
