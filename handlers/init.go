package forum

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	CreateUsersTale := `CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(CreateUsersTale)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	if !tableExists(db, "categories") {
		CreateCategoriestable := `CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		categorie TEXT NOT NULL
	);`
		_, err = db.Exec(CreateCategoriestable)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
		WriteCategories()
	}

	CreatePostsTable := `CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY AUTOINCREMENT ,
			title TEXT ,
			content TEXT ,
			user_id INTEGER ,
			category_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
	);`
	_, err = db.Exec(CreatePostsTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("TABLE successfully added")
}
