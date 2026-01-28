package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) CreateCategory(category models.Category) (models.Category, error) {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at"
	row := r.db.QueryRow(query, category.Name, category.Description)
	var newCategory models.Category
	err := row.Scan(&newCategory.ID, &newCategory.Name, &newCategory.Description, &newCategory.CreatedAt)

	if err != nil {
		return models.Category{}, fmt.Errorf("failed to create category: %w", err)
	}
	return newCategory, nil
}

func (r *CategoryRepository) GetCategoryByID(id string) (models.Category, error) {
	query := "SELECT id, name, description, created_at FROM categories WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var category models.Category
	err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, nil
		}
		return models.Category{}, fmt.Errorf("failed to get category by id %s : %w", id, err)
	}
	return category, nil
}

func (r *CategoryRepository) UpdateCategoryByID(id string, category models.Category) (models.Category, error) {
	query := "UPDATE categories SET name = $2, description = $3 WHERE id = $1 RETURNING id, name, description, created_at"
	row := r.db.QueryRow(query, id, category.Name, category.Description)
	var updatedCategory models.Category
	err := row.Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description, &updatedCategory.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, nil
		}
		return models.Category{}, fmt.Errorf("failed to update category by id %s : %w", id, err)
	}
	return updatedCategory, nil
}

func (r *CategoryRepository) DeleteCategoryByID(id string) (models.Category, error) {
	query := "DELETE FROM categories WHERE id = $1 RETURNING id, name, description, created_at"
	row := r.db.QueryRow(query, id)

	var deletedCategory models.Category
	err := row.Scan(&deletedCategory.ID, &deletedCategory.Name, &deletedCategory.Description, &deletedCategory.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, nil
		}
		return models.Category{}, fmt.Errorf("failed to delete category by id %s : %w", id, err)
	}

	return deletedCategory, nil
}
