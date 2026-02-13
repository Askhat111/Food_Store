package handlers

import (
	"log"
	"time"

	"foodstoreteam/internal/db"
)

func StartAlertWorker() {
	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		checkLowStock()
		checkExpiringProducts()
	}
}

func checkLowStock() {
	rows, err := db.DB.Query(`
        SELECT name, quantity, min_stock_threshold
        FROM products 
        WHERE quantity < min_stock_threshold`)

	if err != nil {
		log.Println("Error checking low stock:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var name string
		var qty, min_stock_threshold int
		if err := rows.Scan(&name, &qty, &min_stock_threshold); err == nil {
			log.Printf("LOW STOCK: %s - Only %d left (minimum: %d)", name, qty, min_stock_threshold)
			count++
		}
	}

	if count > 0 {
		log.Printf("Total low stock items: %d\n", count)
	}
}

func checkExpiringProducts() {
	rows, err := db.DB.Query(`
        SELECT name, expiry_date, quantity 
        FROM products 
        WHERE expiry_date IS NOT NULL 
        AND expiry_date != ''
        AND date(expiry_date) <= date('now', '+7 days')
        ORDER BY expiry_date ASC`)

	if err != nil {
		log.Println("Error checking expiring products:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var name, expiryDate string
		var qty int
		if err := rows.Scan(&name, &expiryDate, &qty); err == nil {
			expiryTime, _ := time.Parse("2006-01-02", expiryDate)
			daysLeft := int(time.Until(expiryTime).Hours() / 24)

			if daysLeft < 0 {
				log.Printf("EXPIRED: %s - Expired %d days ago (Qty: %d)", name, -daysLeft, qty)
			} else if daysLeft == 0 {
				log.Printf("EXPIRES TODAY: %s (Qty: %d)", name, qty)
			} else if daysLeft <= 3 {
				log.Printf("EXPIRES SOON: %s - %d days left (Qty: %d)", name, daysLeft, qty)
			} else {
				log.Printf("EXPIRING: %s - %d days left (Qty: %d)", name, daysLeft, qty)
			}
			count++
		}
	}

	if count > 0 {
		log.Printf("Total expiring/expired items: %d\n", count)
	}
}
