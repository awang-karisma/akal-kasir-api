package models

type Product struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Stock      int        `json:"stock"`
	CreatedAt  string     `json:"created_at"`
	Categories []Category `json:"categories"`
}

type ProductWithoutCategories struct {
	Product
	Categories []Category `json:"categories,omitempty"`
}
