package handlers

import (
	"database/sql"
	"log"
	"time"

	"foodstoreteam/internal/db"
)

func StartLowStockWorker() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		err := db.WithRead(func(sqlDB *sql.DB) error {
			rows, err := sqlDB.Query(
				`SELECT name, quantity, min_stock_threshold FROM products WHERE quantity < min_stock_threshold`)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var name string
				var qty, min int
				if err := rows.Scan(&name, &qty, &min); err != nil {
					return err
				}
				log.Printf("[LOW STOCK] %s: %d (threshold %d)\n", name, qty, min)
			}
			return rows.Err()
		})
		if err != nil {
			log.Println("low-stock worker error:", err)
		}
	}
}
