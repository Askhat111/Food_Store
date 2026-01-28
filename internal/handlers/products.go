package handlers

import (
	"fmt"
	"foodstoreteam/internal/db"
)

func AddProduct() {
	fmt.Print("Name: ")
	fmt.Println("AddProduct skeleton - Merey")
}

func ListProducts() {
	rows, _ := db.DB.Query("SELECT * FROM products")
	defer rows.Close()
	fmt.Println("ListProducts - Merey")
}
