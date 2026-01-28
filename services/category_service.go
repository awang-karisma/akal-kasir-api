package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	return s.repo.GetCategories()
}

func (s *CategoryService) CreateCategory(category models.Category) (models.Category, error) {
	return s.repo.CreateCategory(category)
}

func (s *CategoryService) GetCategoryByID(id string) (models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) UpdateCategoryByID(id string, category models.Category) (models.Category, error) {
	return s.repo.UpdateCategoryByID(id, category)
}

func (s *CategoryService) DeleteCategoryByID(id string) (models.Category, error) {
	return s.repo.DeleteCategoryByID(id)
}
