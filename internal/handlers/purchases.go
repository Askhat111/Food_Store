package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

func HandlePurchases(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		RequireRole("owner")(createPurchase)(w, r)
	case http.MethodGet:
		AuthMiddleware(listPurchases)(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createPurchase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Supplier  string  `json:"supplier"`
		Price     float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 || req.Price <= 0 {
		http.Error(w, "Invalid quantity or price", http.StatusBadRequest)
		return
	}

	//start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	//checking if product exists
	var productName string
	err = tx.QueryRow("SELECT name FROM products WHERE id = ?", req.ProductID).Scan(&productName)
	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//updating
	_, err = tx.Exec(`
        UPDATE products 
        SET quantity = quantity + ?, price = ?
        WHERE id = ?`, req.Quantity, req.Price, req.ProductID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//purchase record
	now := time.Now().Format(time.RFC3339)
	result, err := tx.Exec(`
        INSERT INTO purchases (product_id, quantity, purchase_date, supplier, price)
        VALUES (?, ?, ?, ?, ?)`,
		req.ProductID, req.Quantity, now, req.Supplier, req.Price)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//commiting transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	purchaseID, _ := result.LastInsertId()
	total := req.Price * float64(req.Quantity)

	purchase := models.Purchase{
		ID:           int(purchaseID),
		ProductID:    req.ProductID,
		ProductName:  productName,
		Quantity:     req.Quantity,
		PurchaseDate: now,
		Supplier:     req.Supplier,
		Price:        req.Price,
		Total:        total,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(purchase)
}

func listPurchases(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
        SELECT p.id, p.product_id, p.quantity, p.purchase_date, p.supplier, p.price,
               pr.name as product_name
        FROM purchases p
        JOIN products pr ON p.product_id = pr.id
        ORDER BY p.purchase_date DESC
        LIMIT 100`)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	purchases := []models.Purchase{}
	for rows.Next() {
		var p models.Purchase
		if err := rows.Scan(&p.ID, &p.ProductID, &p.Quantity, &p.PurchaseDate,
			&p.Supplier, &p.Price, &p.ProductName); err != nil {
			continue
		}
		p.Total = p.Price * float64(p.Quantity)
		purchases = append(purchases, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchases)
}
