package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTransactions() ([]models.Transaction, error) {
	return s.repo.GetTransactions()
}

func (s *TransactionService) GetTransactionsRange(from string, to string) ([]models.Transaction, error) {
	return s.repo.GetTransactionsRange(from, to)
}
