package db

import (
	"github.com/go-pg/pg/v9"
	"log"
	"os"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     getEnv("QOVERY_DATABASE_MY_DB_USERNAME", "postgres"),
		Password: getEnv("QOVERY_DATABASE_MY_DB_PASSWORD", "postgres"),
		Addr:     getEnv("QOVERY_DATABASE_MY_DB_HOST", "localhost") + ":" + getEnv("QOVERY_DATABASE_MY_DB_PORT", "5432"),
		Database: getEnv("QOVERY_DATABASE_MY_DB_USERNAME", "postgres"),
	}

	var db *pg.DB = pg.Connect(opts)

	if db == nil {
		log.Printf("Failed to connect to PotsgreSQL")
		os.Exit(100)
	}

	log.Printf("Connected to PostgreSQL")

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
