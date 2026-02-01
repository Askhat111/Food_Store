package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

func LowStockReport(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	err := db.WithRead(func(sqlDB *sql.DB) error {
		rows, err := sqlDB.Query(
			`SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold
			 FROM products WHERE quantity < min_stock_threshold`)
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
