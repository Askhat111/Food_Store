package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

// products list and creating
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listProducts(w, r)
	case http.MethodPost:
		RequireRole("owner")(createProduct)(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProduct(w, r, id)
	case http.MethodPut:
		RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
			updateProduct(w, r, id)
		})(w, r)
	case http.MethodDelete:
		RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
			deleteProduct(w, r, id)
		})(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// new product
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if p.Name == "" || p.Price <= 0 || p.Quantity < 0 {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	if p.MinStockThreshold == 0 {
		p.MinStockThreshold = 5
	}

	result, err := db.DB.Exec(`
        INSERT INTO products (name, price, quantity, category, expiry_date, min_stock_threshold)
        VALUES (?, ?, ?, ?, ?, ?)`,
		p.Name, p.Price, p.Quantity, p.Category, p.ExpiryDate, p.MinStockThreshold)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")

	query := `SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold FROM products WHERE 1=1`
	args := []interface{}{}

	//search filter
	if search != "" {
		query += " AND (name LIKE ? OR category LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	//category filter
	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category, &p.ExpiryDate, &p.MinStockThreshold); err != nil {
			continue
		}

		//days to expiry
		if p.ExpiryDate != "" {
			if expiryTime, err := time.Parse("2006-01-02", p.ExpiryDate); err == nil {
				p.DaysToExpiry = int(time.Until(expiryTime).Hours() / 24)
			}
		}

		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request, id int) {
	var p models.Product
	err := db.DB.QueryRow(`
        SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold 
        FROM products WHERE id = ?`, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category, &p.ExpiryDate, &p.MinStockThreshold)

	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//days to expiry
	if p.ExpiryDate != "" {
		if expiryTime, err := time.Parse("2006-01-02", p.ExpiryDate); err == nil {
			p.DaysToExpiry = int(time.Until(expiryTime).Hours() / 24)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if p.Name == "" || p.Price <= 0 || p.Quantity < 0 {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
        UPDATE products 
        SET name = ?, price = ?, quantity = ?, category = ?, expiry_date = ?, min_stock_threshold = ?
        WHERE id = ?`,
		p.Name, p.Price, p.Quantity, p.Category, p.ExpiryDate, p.MinStockThreshold, id)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM sales WHERE product_id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Cannot delete product with sales history", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
