package db

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	DB  *sql.DB
	mux sync.RWMutex
)

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./data/foodstore.db")
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			quantity INTEGER NOT NULL,
			category TEXT,
			expiry_date TEXT,
			min_stock_threshold INTEGER DEFAULT 5
		);
	`)
	return err
}

func WithRead(f func(*sql.DB) error) error {
	mux.RLock()
	defer mux.RUnlock()
	return f(DB)
}

func WithWrite(f func(*sql.DB) error) error {
	mux.Lock()
	defer mux.Unlock()
	return f(DB)
}
