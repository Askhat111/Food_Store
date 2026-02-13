package db

import (
	"database/sql"
	"foodstoreteam/internal/models"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "foodstore.db")
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
        
        CREATE TABLE IF NOT EXISTS sales (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            product_id INTEGER NOT NULL,
            quantity INTEGER NOT NULL,
            sale_date TEXT NOT NULL,
            total REAL NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY(product_id) REFERENCES products(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
        
        CREATE TABLE IF NOT EXISTS purchases (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            product_id INTEGER NOT NULL,
            quantity INTEGER NOT NULL,
            purchase_date TEXT NOT NULL,
            supplier TEXT,
            price REAL NOT NULL,
            FOREIGN KEY(product_id) REFERENCES products(id)
        );
        
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            name TEXT NOT NULL,
            role TEXT NOT NULL,
            password TEXT NOT NULL
        );
        
        INSERT OR IGNORE INTO users (id, username, name, role, password) VALUES 
        (1, 'admin', 'Store Owner', 'owner', 'admin123'),
        (2, 'cashier', 'Cashier', 'staff', 'cash123');
    `)
	return err
}

func Authenticate(username, password string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(`
		SELECT id, username, name, role 
		FROM users WHERE username = ? AND password = ?`,
		username, password).Scan(&user.ID, &user.Username, &user.Name, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(`
		SELECT id, username, name, role 
		FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Username, &user.Name, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
