package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetAllProducts() ([]models.Product, error) {
	args := m.Called()

	var products []models.Product
	if args.Get(0) != nil {
		products = args.Get(0).([]models.Product)
	}
	return products, args.Error(1)
}
