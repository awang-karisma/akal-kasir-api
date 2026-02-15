package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalAmount int64 = 0
	details := make([]models.TransactionDetail, 0)
	for _, item := range items {
		var productDetail models.ProductDetail
		err = tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productDetail.ID, &productDetail.Name, &productDetail.Price, &productDetail.Stock)
		if err != nil {
			return nil, err
		}
		if productDetail.Stock < item.Quantity {
			return nil, fmt.Errorf("Insufficient stock for item %s, stock is %d but requested %d", item.ProductID, productDetail.Stock, item.Quantity)
		}

		subTotal := productDetail.Price * int64(item.Quantity)
		totalAmount += subTotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   productDetail.ID,
			ProductName: productDetail.Name,
			Quantity:    item.Quantity,
			Subtotal:    subTotal,
			Price:       productDetail.Price,
		})

	}

	var transactionId string
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionId)
	if err != nil {
		return nil, err
	}

	for _, detail := range details {
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", transactionId, detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionId,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
