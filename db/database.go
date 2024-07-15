package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Fatal("POSTGRES_URI environment variable is not set")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS skills (
			key VARCHAR(100) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			logo TEXT,
			tags TEXT[]
	);
	`
	_, err = DB.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	fmt.Println("Successfully connected to the database!")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
