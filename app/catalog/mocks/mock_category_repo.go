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
