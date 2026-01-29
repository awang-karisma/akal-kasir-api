package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kasir-api/models"
	"strings"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts() ([]models.Product, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.created_at,
			COALESCE(
				json_agg(json_build_object(
					'id', c.id,
					'name', c.name,
					'description', c.description,
					'created_at', c.created_at
				)) FILTER (WHERE c.id IS NOT NULL),
				'[]'::json
			) AS categories
		FROM products p
		LEFT JOIN product_categories pc ON p.id = pc.product_id
		LEFT JOIN categories c ON pc.category_id = c.id
		GROUP BY p.id, p.name;
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		var categories string
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &categories)
		if err != nil {
			return nil, err
		}
		json.NewDecoder(strings.NewReader(categories)).Decode(&product.Categories)
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(product models.Product) (models.Product, error) {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id, name, price, stock"
	row := r.db.QueryRow(query, product.Name, product.Price, product.Stock)
	var newProduct models.Product
	err := row.Scan(&newProduct.ID, &newProduct.Name, &newProduct.Price, &newProduct.Stock)

	if err != nil {
		return models.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return newProduct, nil
}

func (r *ProductRepository) GetProductByID(id string) (models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.created_at,
		       c.id, c.name, c.description, c.created_at
		FROM products p
		LEFT JOIN product_categories pc ON p.id = pc.product_id
		LEFT JOIN categories c ON pc.category_id = c.id
		WHERE p.id = $1
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to get product by id %s : %w", id, err)
	}
	defer rows.Close()

	var product models.Product
	product.Categories = []models.Category{}
	for rows.Next() {
		var category models.Category
		var categoryID, categoryName, categoryDescription, categoryCreatedAt *string

		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt,
			&categoryID, &categoryName, &categoryDescription, &categoryCreatedAt)
		if err != nil {
			return models.Product{}, fmt.Errorf("failed to scan product: %w", err)
		}

		if categoryID != nil {
			category = models.Category{
				ID:          *categoryID,
				Name:        *categoryName,
				Description: *categoryDescription,
				CreatedAt:   *categoryCreatedAt,
			}
			product.Categories = append(product.Categories, category)
		}
	}

	if product.ID == "" {
		return models.Product{}, nil
	}
	return product, nil
}

func (r *ProductRepository) UpdateProductByID(id string, product models.Product) (models.Product, error) {
	query := "UPDATE products SET name = $2, price = $3, stock = $4 WHERE id = $1 RETURNING id, name, price, stock"
	row := r.db.QueryRow(query, id, product.Name, product.Price, product.Stock)
	var updatedProduct models.Product
	err := row.Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Price, &updatedProduct.Stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, nil
		}
		return models.Product{}, fmt.Errorf("failed to update product by id %s : %w", id, err)
	}
	return updatedProduct, nil
}

func (r *ProductRepository) DeleteProductByID(id string) (models.Product, error) {
	query := "DELETE FROM products WHERE id = $1 RETURNING id, name, price, stock"
	row := r.db.QueryRow(query, id)

	var deletedProduct models.Product
	err := row.Scan(&deletedProduct.ID, &deletedProduct.Name, &deletedProduct.Price, &deletedProduct.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, nil
		}
		return models.Product{}, fmt.Errorf("failed to delete product by id %s : %w", id, err)
	}

	return deletedProduct, nil
}

func (r *ProductRepository) AddCategoryToProduct(productID, categoryID string) error {
	query := "INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, productID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to add category to product: %w", err)
	}
	return nil
}

func (r *ProductRepository) RemoveCategoryFromProduct(productID, categoryID string) error {
	query := "DELETE FROM product_categories WHERE product_id = $1 AND category_id = $2"
	_, err := r.db.Exec(query, productID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove category from product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetCategoriesByProductID(productID string) ([]models.Category, error) {
	query := `
		SELECT c.id, c.name, c.description, c.created_at
		FROM categories c
		INNER JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
		ORDER BY c.name
	`
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by product id %s : %w", productID, err)
	}
	defer rows.Close()

	// var categories []models.Category
	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
