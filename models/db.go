package models

import (
	"database/sql"
	"log"
)

var db *sql.DB

func ConnectDB() {
	var err error
	db, err = sql.Open("sqlite", "./urls.db")
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table", err)
	}
}
