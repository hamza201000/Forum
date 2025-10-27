package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	createTableSQL := `CREATE TABLE IF NOT EXIST users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NUL,
			email TEXT NOT NUL UNIQUE
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
