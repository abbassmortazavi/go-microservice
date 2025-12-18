package database

import (
	"abbassmortazavi/go-microservice/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connect() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Name)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxConn)
	duration, err := time.ParseDuration(cfg.MaxIdleTimeout)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	DB = db
}
