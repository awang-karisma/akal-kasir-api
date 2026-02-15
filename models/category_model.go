package models

import "time"

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CategoryWithProducts struct {
	Category
	Products []ProductWithoutCategories `json:"products"`
}
