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
	// Define methods for category operations
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// / Helpers
func (r *CategoryRepository) handleGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("category not found: %w", ErrCategoryNotFound)
	}

	// Map other GORM errors to your sentinel errors
	return fmt.Errorf("database error: %w", ErrNoProductAssociated)
}
