package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(db *sql.DB) {
	dbURL := os.Getenv("DB_ADDRESS")
	if dbURL == "" {
		dbURL = "postgres://root:root@localhost:5432/microservice_db?sslmode=disable"
	}

	// ایجاد migration instance
	m, err := migrate.New(
		"file://services/auth-service/migrations", // مسیر فایل‌های migration
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}
	defer m.Close()

	// اجرای migration
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to run migrations: %s", err)
	}

	log.Println("Migrations applied successfully")

}
