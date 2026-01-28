package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
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
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, nil
		}
		return models.Product{}, fmt.Errorf("failed to get product by id %s : %w", id, err)
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
