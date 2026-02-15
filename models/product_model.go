package models

import "time"

type Product struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Price      int64      `json:"price"`
	Stock      int        `json:"stock"`
	CreatedAt  time.Time  `json:"created_at"`
	Categories []Category `json:"categories"`
}

type ProductWithoutCategories struct {
	Product
	Categories []Category `json:"categories,omitempty"`
}
