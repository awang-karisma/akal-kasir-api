package models

type Transaction struct {
	ID          string              `json:"id"`
	TotalAmount int64               `json:"total_amount"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int64  `json:"subtotal"`
	Price         int64  `json:"price"`
}

// NOTE : Move this struct to product model if it is used in multiple places
type ProductDetail struct {
	ID    string `json:"product_id"`
	Name  string `json:"product_name"`
	Price int64  `json:"product_price"`
	Stock int    `json:"product_stock"`
}

type CheckoutItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
