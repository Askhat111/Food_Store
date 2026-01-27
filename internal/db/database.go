package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "./data/foodstore.db")
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			price REAL,
			quantity INTEGER,
			category TEXT,
			expiry TEXT
		)
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("SQLite initialized")
}
