package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrNoProductAssociated = errors.New("no product associated with this category")
)

type CategoryInterface interface {
	GetAllCategories() ([]Category, error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, r.handleGormError(err)
	}
	return categories, nil
}

// / Helpers
func (r *CategoryRepository) handleGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("category not found: %w", ErrCategoryNotFound)
	}

	// Map other GORM errors to your sentinel errors
	return fmt.Errorf("database error: %w", ErrNoProductAssociated)
}
