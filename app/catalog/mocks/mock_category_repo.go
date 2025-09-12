package mocks

import "github.com/mytheresa/go-hiring-challenge/models"

func (m *MockCatalogRepository) GetAllCategories() ([]models.Category, error) {
	args := m.Called()

	var categories []models.Category
	if args.Get(0) != nil {
		categories = args.Get(0).([]models.Category)
	}
	return categories, args.Error(1)
}

func (m *MockCatalogRepository) CreateCategory(code, name string) (*models.Category, error) {
	args := m.Called(code, name)

	var category *models.Category
	if args.Get(0) != nil {
		category = args.Get(0).(*models.Category)
	}
	return category, args.Error(1)
}
