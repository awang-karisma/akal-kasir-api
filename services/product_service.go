package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repo.GetProducts()
}

func (s *ProductService) CreateProduct(product models.Product) (models.Product, error) {
	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id string) (models.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) UpdateProductByID(id string, product models.Product) (models.Product, error) {
	return s.repo.UpdateProductByID(id, product)
}

func (s *ProductService) DeleteProductByID(id string) (models.Product, error) {
	return s.repo.DeleteProductByID(id)
}
