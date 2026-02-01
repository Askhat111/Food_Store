package models

type Product struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	Category          string  `json:"category"`
	ExpiryDate        string  `json:"expiry_date"`
	MinStockThreshold int     `json:"min_stock_threshold"`
}
