package catalog

import (
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// Product represents a product in API responses
type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category,omitempty"`
}

// Category represents a category in API responses
type Category struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ProductsResponse represents the response for the /catalog endpoint
type ProductsResponse struct {
	Products []Product `json:"products"`
}

// CategoriesResponse represents the response for the /categories endpoint
type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type CatalogRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetAllCategories() ([]models.Category, error)
}

type CatalogService struct {
	ProductsRepo models.ProductInterface
	CategoryRepo models.CategoryInterface
}

type CatalogHandler struct {
	repo CatalogRepository
}

func NewCatalogService(pr models.ProductInterface, cr models.CategoryInterface) CatalogRepository {
	return &CatalogService{
		ProductsRepo: pr,
		CategoryRepo: cr,
	}
}

func NewCatalogHandler(r CatalogRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleRetrieveProducts(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllProducts()
	if err != nil {
		// Handle different types of errors with proper HTTP status codes
		h.handleError(w, err)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				ID:   p.Category.ID,
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	// Return the products as a JSON response
	response := ProductsResponse{
		Products: products,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleRetrieveCategories(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllCategories()
	if err != nil {
		// Handle different types of errors with proper HTTP status codes
		h.handleError(w, err)
		return
	}

	// Map response
	categories := make([]Category, len(res))
	for i, c := range res {
		categories[i] = Category{
			ID:   c.ID,
			Code: c.Code,
			Name: c.Name,
		}
	}

	// Return the products as a JSON response
	response := CategoriesResponse{
		Categories: categories,
	}

	api.OKResponse(w, response)
}

func (s *CatalogService) GetAllProducts() ([]models.Product, error) {
	return s.ProductsRepo.GetAllProducts()
}

func (s *CatalogService) GetAllCategories() ([]models.Category, error) {
	return s.CategoryRepo.GetAllCategories()
}

// handleError checks for specific error types and returns appropriate HTTP status codes
func (h *CatalogHandler) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrProductNotFound) || errors.Is(err, models.ErrCategoryNotFound) {
		api.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if errors.Is(err, models.ErrDatabaseConnection) {
		api.ErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	// For all other errors
	api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
}
