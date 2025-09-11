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

type ProductQueryParams struct {
	Limit        int
	Offset       int
	CategoryName string
	PriceMax     float64
}

type ProductInterface interface {
	GetAllProducts(params ProductQueryParams) ([]Product, error)
	GetProductByCode(code string) (*Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(params ProductQueryParams) ([]Product, error) {
	var products []Product
	query := r.buildQuery(params)

	if err := query.Limit(params.Limit).Offset(params.Offset).Find(&products).Error; err != nil {
		return nil, r.handleGormError(err)
	}
	return products, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Variants").Preload("Category").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, r.handleGormError(err)
	}
	return &product, nil
}

// / Helpers
func (r *ProductsRepository) handleGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("product not found: %w", ErrProductNotFound)
	}

	// Map other GORM errors to your sentinel errors
	return fmt.Errorf("database error: %w", ErrDatabaseConnection)
}

func (r *ProductsRepository) buildQuery(params ProductQueryParams) *gorm.DB {
	query := r.db.Model(&Product{}).Preload("Variants").Preload("Category")
	query = applyCategoryFilter(query, params.CategoryName)
	query = applyPriceFilter(query, params.PriceMax)
	return query
}

func applyCategoryFilter(query *gorm.DB, categoryName string) *gorm.DB {
	if categoryName != "" {
		return query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name = ?", categoryName)
	}
	return query
}

func applyPriceFilter(query *gorm.DB, priceMax float64) *gorm.DB {
	if priceMax > 0 {
		return query.Where("price <= ?", priceMax)
	}
	return query
}
