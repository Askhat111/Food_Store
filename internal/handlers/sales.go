package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

func HandleSales(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AuthMiddleware(createSale)(w, r)
	case http.MethodGet:
		AuthMiddleware(listSales)(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createSale(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*models.User)

	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be positive", http.StatusBadRequest)
		return
	}

	//starting transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	//product details
	var productName string
	var price float64
	var currentQty int
	err = tx.QueryRow("SELECT name, price, quantity FROM products WHERE id = ?", req.ProductID).
		Scan(&productName, &price, &currentQty)

	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//stock
	if currentQty < req.Quantity {
		http.Error(w, "Insufficient stock", http.StatusBadRequest)
		return
	}

	//new product quantity(+)
	_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", req.Quantity, req.ProductID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//sale record
	total := price * float64(req.Quantity)
	now := time.Now().Format(time.RFC3339)

	result, err := tx.Exec(`
        INSERT INTO sales (product_id, quantity, sale_date, total, user_id)
        VALUES (?, ?, ?, ?, ?)`,
		req.ProductID, req.Quantity, now, total, user.ID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//commiting transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	saleID, _ := result.LastInsertId()
	sale := models.Sale{
		ID:          int(saleID),
		ProductID:   req.ProductID,
		ProductName: productName,
		Quantity:    req.Quantity,
		SaleDate:    now,
		Total:       total,
		UserID:      user.ID,
		UserName:    user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sale)
}

func listSales(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*models.User)
	date := r.URL.Query().Get("date")

	query := `
        SELECT s.id, s.product_id, s.quantity, s.sale_date, s.total, s.user_id,
               p.name as product_name, u.name as user_name
        FROM sales s
        JOIN products p ON s.product_id = p.id
        JOIN users u ON s.user_id = u.id
        WHERE 1=1`

	args := []interface{}{}

	if user.Role == "staff" {
		query += " AND s.user_id = ?"
		args = append(args, user.ID)
	}

	//filter by date
	if date != "" {
		query += " AND DATE(s.sale_date) = ?"
		args = append(args, date)
	}

	query += " ORDER BY s.sale_date DESC LIMIT 100"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	sales := []models.Sale{}
	for rows.Next() {
		var s models.Sale
		if err := rows.Scan(&s.ID, &s.ProductID, &s.Quantity, &s.SaleDate, &s.Total,
			&s.UserID, &s.ProductName, &s.UserName); err != nil {
			continue
		}
		sales = append(sales, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}
