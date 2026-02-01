package main

import (
	"log"
	"net/http"

	"foodstoreteam/internal/db"
	"foodstoreteam/internal/handlers"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	go handlers.StartLowStockWorker()

	mux := http.NewServeMux()

	mux.HandleFunc("/products", handlers.HandleProducts)
	mux.HandleFunc("/products/", handlers.HandleProductByID)
	mux.HandleFunc("/reports/low-stock", handlers.LowStockReport)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Food Store HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
