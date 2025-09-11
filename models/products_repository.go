package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrDatabaseConnection = errors.New("database connection error")
)

type ProductInterface interface {
	GetAllProducts() ([]Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Variants").Preload("Category").Find(&products).Error; err != nil {
		return nil, r.handleGormError(err)
	}
	return products, nil
}

// / Helpers
func (r *ProductsRepository) handleGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("product not found: %w", ErrProductNotFound)
	}

	// Map other GORM errors to your sentinel errors
	return fmt.Errorf("database error: %w", ErrDatabaseConnection)
}
