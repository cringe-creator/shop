package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE is not set")
	}

	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return &Database{DB: conn}
}
