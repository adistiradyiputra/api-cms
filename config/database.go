package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Set default port if empty
	port := ENV.DBPort
	if port == "" {
		port = "5432" // Default PostgreSQL port
	}

	// Check if required env vars are set
	if ENV.DBHost == "" || ENV.DBUser == "" || ENV.DBPassword == "" || ENV.DBName == "" {
		log.Println("Warning: Database environment variables not fully set")
		log.Printf("Host: %s, User: %s, DB: %s, Port: %s", ENV.DBHost, ENV.DBUser, ENV.DBName, port)
		return
	}

	// Ambil data dari ENV struct
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		ENV.DBHost, ENV.DBUser, ENV.DBPassword, ENV.DBName, port,
	)

	log.Printf("Connecting to database with DSN: host=%s user=%s dbname=%s port=%s", ENV.DBHost, ENV.DBUser, ENV.DBName, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Gagal konek ke database: %v", err)
		log.Printf("DSN: %s", dsn)
		return
	}

	log.Println("Database connected successfully")
}
