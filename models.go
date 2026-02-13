package models

type Product struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	Category          string  `json:"category"`
	ExpiryDate        string  `json:"expiry_date"`
	MinStockThreshold int     `json:"min_stock_threshold"`
	DaysToExpiry      int     `json:"days_to_expiry,omitempty"`
}

type Sale struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	SaleDate    string  `json:"sale_date"`
	Total       float64 `json:"total"`
	UserID      int     `json:"user_id"`
	UserName    string  `json:"user_name"`
}

type Purchase struct {
	ID           int     `json:"id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Quantity     int     `json:"quantity"`
	PurchaseDate string  `json:"purchase_date"`
	Supplier     string  `json:"supplier"`
	Price        float64 `json:"price"`
	Total        float64 `json:"total"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type Alert struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Product string `json:"product"`
	Value   int    `json:"value"`
}

type DailySales struct {
	Date    string  `json:"date"`
	Total   float64 `json:"total"`
	Count   int     `json:"count"`
	AvgSale float64 `json:"avg_sale"`
}
