package main

import (
	"fmt"
	"foodstoreteam/internal/db"
	"foodstoreteam/internal/handlers"
)

func main() {
	db.InitDB()
	fmt.Println("1. Add Product")
	fmt.Println("2. List Products")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		handlers.AddProduct()
	case "2":
		handlers.ListProducts()
	}
}
