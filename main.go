package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Подключаемся к БД
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url(
		    id SERIAL PRIMARY KEY,
		    alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
`)
	if err != nil {
		log.Fatal(err)
	}

	var insertedID int
	err = db.QueryRow(
		context.Background(),
		`INSERT INTO url (alias, url) VALUES ($1, $2),
        	 ON CONFLICT (alias) DO UPDATE SET url = EXCLUDED.url
        	 RETURNING id;`,
		"example12",
		"https://example.com",
	).Scan(&insertedID)

	if err != nil {
		log.Fatalf("insert failed: %v\n", err)
	}

	fmt.Printf("Операция успешно выполнена ID: %d\n", insertedID)
}
