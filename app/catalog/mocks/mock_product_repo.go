package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type MockCatalogRepository struct {
	mock.Mock
}

func (m *MockCatalogRepository) GetAllProducts(params models.ProductQueryParams) ([]models.Product, error) {
	args := m.Called(params)

	var products []models.Product
	if args.Get(0) != nil {
		products = args.Get(0).([]models.Product)
	}
	return products, args.Error(1)
}
