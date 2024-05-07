package db

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDB() (*sql.DB, error) {
	connStr := "postgres://postgres:ryd12@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Connected to the database!")

	return db, nil
}
