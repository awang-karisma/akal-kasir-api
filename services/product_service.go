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

func (s *ProductService) GetProducts(name string) ([]models.Product, error) {
	return s.repo.GetProducts(name)
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

func (s *ProductService) AddCategoryToProduct(productID, categoryID string) error {
	return s.repo.AddCategoryToProduct(productID, categoryID)
}

func (s *ProductService) RemoveCategoryFromProduct(productID, categoryID string) error {
	return s.repo.RemoveCategoryFromProduct(productID, categoryID)
}

func (s *ProductService) GetCategoriesByProductID(productID string) ([]models.Category, error) {
	return s.repo.GetCategoriesByProductID(productID)
}
