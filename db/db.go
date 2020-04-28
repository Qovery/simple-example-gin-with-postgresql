package db

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"log"
	"os"
	"time"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     getEnv("QOVERY_DATABASE_MY_DB_USERNAME", "postgres"),
		Password: getEnv("QOVERY_DATABASE_MY_DB_PASSWORD", "postgres"),
		Addr:     getEnv("QOVERY_DATABASE_MY_DB_HOST", "localhost") + ":" + getEnv("QOVERY_DATABASE_MY_DB_PORT", "5432"),
		Database: getEnv("QOVERY_DATABASE_MY_DB_DATABASE", "postgres"),
		OnConnect: func(conn *pg.Conn) error {
			return CreateJokeTable(conn)
		},
	}

	var db *pg.DB = pg.Connect(opts)

	if db == nil {
		log.Printf("Failed to connect to PotsgreSQL")
		os.Exit(100)
	}

	log.Printf("Connected to PostgreSQL")

	return db
}

func CreateJokeTable(conn *pg.Conn) error {
	opts := orm.CreateTableOptions{
		IfNotExists: true,
	}

	createError := conn.CreateTable(&Joke{}, &opts)

	if createError != nil {
		log.Printf("Error while creating Jokes table, reason: %v\n", createError)
		return createError
	}

	log.Printf("Jokes table created")
	return nil
}

type Joke struct {
	Content string    `json:"content"`
	AddedAt time.Time `json:"added_at"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
