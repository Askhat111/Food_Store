package main

import (
	"encoding/json"
	"log"
	"net/http"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/handlers"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal("Database initialization error:", err)
	}

	go handlers.StartAlertWorker()

	mux := http.NewServeMux()

	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	mux.HandleFunc("/api/login", cors(handlers.HandleLogin))
	mux.HandleFunc("/api/products", cors(handlers.HandleProducts))
	mux.HandleFunc("/api/products/", cors(handlers.HandleProductByID))
	mux.HandleFunc("/api/sales", cors(handlers.HandleSales))
	mux.HandleFunc("/api/purchases", cors(handlers.HandlePurchases))
	mux.HandleFunc("/api/reports/low-stock", cors(handlers.LowStockReport))
	mux.HandleFunc("/api/reports/expiring", cors(handlers.ExpiringProductsReport))
	mux.HandleFunc("/api/reports/daily-sales", cors(handlers.DailySalesReport))
	mux.HandleFunc("/api/alerts", cors(handlers.GetAlerts))
	mux.HandleFunc("/api/health", cors(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
		} else {
			http.ServeFile(w, r, "web"+r.URL.Path)
		}
	})

	port := ":8080"
	log.Printf("Food Store API running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
