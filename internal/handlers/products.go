package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

// products - GET (list), POST (create)
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// products/{id} - PUT, DELETE
func HandleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if err := db.WithWrite(func(sqlDB *sql.DB) error {
		res, err := sqlDB.Exec(
			`INSERT INTO products (name, price, quantity, category, expiry_date, min_stock_threshold)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			p.Name, p.Price, p.Quantity, p.Category, p.ExpiryDate, p.MinStockThreshold,
		)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		p.ID = int(id)
		return nil
	}); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	err := db.WithRead(func(sqlDB *sql.DB) error {
		rows, err := sqlDB.Query(
			`SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold FROM products`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity,
				&p.Category, &p.ExpiryDate, &p.MinStockThreshold); err != nil {
				return err
			}
			products = append(products, p)
		}
		return rows.Err()
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}

func updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	p.ID = id
	err := db.WithWrite(func(sqlDB *sql.DB) error {
		_, err := sqlDB.Exec(
			`UPDATE products SET name=?, price=?, quantity=?, category=?, expiry_date=?, min_stock_threshold=?
			 WHERE id=?`,
			p.Name, p.Price, p.Quantity, p.Category, p.ExpiryDate, p.MinStockThreshold, p.ID,
		)
		return err
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	err := db.WithWrite(func(sqlDB *sql.DB) error {
		_, err := sqlDB.Exec(`DELETE FROM products WHERE id=?`, id)
		return err
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
