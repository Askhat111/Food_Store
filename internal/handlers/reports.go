package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/models"
)

func LowStockReport(w http.ResponseWriter, r *http.Request) {
	RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.DB.Query(`
            SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold
            FROM products 
            WHERE quantity < min_stock_threshold 
            ORDER BY quantity ASC`)

		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := []models.Product{}
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity,
				&p.Category, &p.ExpiryDate, &p.MinStockThreshold); err != nil {
				continue
			}
			products = append(products, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})(w, r)
}

func ExpiringProductsReport(w http.ResponseWriter, r *http.Request) {
	RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
		days := r.URL.Query().Get("days")
		if days == "" {
			days = "7"
		}

		rows, err := db.DB.Query(`
            SELECT id, name, price, quantity, category, expiry_date, min_stock_threshold
            FROM products 
            WHERE expiry_date IS NOT NULL 
            AND expiry_date != ''
            AND date(expiry_date) <= date('now', '+' || ? || ' days')
            ORDER BY expiry_date ASC`, days)

		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := []models.Product{}
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity,
				&p.Category, &p.ExpiryDate, &p.MinStockThreshold); err != nil {
				continue
			}

			//days to expiry
			if expiryTime, err := time.Parse("2006-01-02", p.ExpiryDate); err == nil {
				p.DaysToExpiry = int(time.Until(expiryTime).Hours() / 24)
			}

			products = append(products, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})(w, r)
}

func DailySalesReport(w http.ResponseWriter, r *http.Request) {
	RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		if startDate == "" {
			startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
		}
		if endDate == "" {
			endDate = time.Now().Format("2006-01-02")
		}

		rows, err := db.DB.Query(`
            SELECT 
                DATE(sale_date) as date,
                SUM(total) as total,
                COUNT(*) as count
            FROM sales
            WHERE DATE(sale_date) BETWEEN ? AND ?
            GROUP BY DATE(sale_date)
            ORDER BY date DESC`, startDate, endDate)

		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		dailySales := []models.DailySales{}
		for rows.Next() {
			var ds models.DailySales
			if err := rows.Scan(&ds.Date, &ds.Total, &ds.Count); err != nil {
				continue
			}
			if ds.Count > 0 {
				ds.AvgSale = ds.Total / float64(ds.Count)
			}
			dailySales = append(dailySales, ds)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dailySales)
	})(w, r)
}

func GetAlerts(w http.ResponseWriter, r *http.Request) {
	RequireRole("owner")(func(w http.ResponseWriter, r *http.Request) {
		alerts := []models.Alert{}

		//low stock alerts
		rows1, err := db.DB.Query(`
            SELECT name, quantity, min_stock_threshold
            FROM products 
            WHERE quantity < min_stock_threshold`)

		if err == nil {
			defer rows1.Close()
			for rows1.Next() {
				var name string
				var qty, minStock int
				if rows1.Scan(&name, &qty, &minStock) == nil {
					alerts = append(alerts, models.Alert{
						Type:    "low_stock",
						Message: "Low stock alert",
						Product: name,
						Value:   qty,
					})
				}
			}
		}

		//expiring products alerts
		rows2, err := db.DB.Query(`
            SELECT name, expiry_date, quantity
            FROM products 
            WHERE expiry_date IS NOT NULL 
            AND expiry_date != ''
            AND date(expiry_date) <= date('now', '+7 days')`)

		if err == nil {
			defer rows2.Close()
			for rows2.Next() {
				var name, expiryDate string
				var qty int
				if rows2.Scan(&name, &expiryDate, &qty) == nil {
					if expiryTime, err := time.Parse("2006-01-02", expiryDate); err == nil {
						daysLeft := int(time.Until(expiryTime).Hours() / 24)
						var message string
						if daysLeft < 0 {
							message = "Expired"
						} else if daysLeft <= 3 {
							message = "Expires soon"
						} else {
							message = "Expiring"
						}

						alerts = append(alerts, models.Alert{
							Type:    "expiring",
							Message: message,
							Product: name,
							Value:   daysLeft,
						})
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alerts)
	})(w, r)
}
