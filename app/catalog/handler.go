package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// Product represents a product in API responses
type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category,omitempty"`
	Variants []Variant `json:"variants,omitempty"`
}

type Variant struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price,omitempty"`
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
	Total    int       `json:"total"`
	Products []Product `json:"products"`
}

// CategoriesResponse represents the response for the /categories endpoint
type CategoriesResponse struct {
	Total      int        `json:"total"`
	Categories []Category `json:"categories"`
}

type CatalogRepository interface {
	GetAllProducts(p models.ProductQueryParams) ([]models.Product, error)
	GetAllCategories() ([]models.Category, error)
	CreateCategory(code, name string) (*models.Category, error)
	GetProductByCode(code string) (*models.Product, error)
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
	queryParams := parseQueryParams(r)
	res, err := h.repo.GetAllProducts(queryParams)
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
		Total:    len(products),
		Products: products,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleRetrieveProductByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Missing product code")
		return
	}

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert variants
	variants := make([]Variant, len(product.Variants))
	for i, v := range product.Variants {
		price := 0.0
		if v.Price.BigInt().Sign() > 0 { // Check if price is not null
			price = v.Price.InexactFloat64()
		}

		variants[i] = Variant{
			ID:    v.ID,
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		}
	}

	// Convert the model product to a response product
	response := Product{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Variants: variants,
		Category: Category{
			ID:   product.Category.ID,
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
	}

	// Return the product as a JSON response
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
		Total:      len(categories),
		Categories: categories,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, err)
		return
	}

	res, err := h.repo.CreateCategory(req.Code, req.Name)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return the newly created category as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Category *models.Category `json:"category"`
	}{
		Category: res,
	}
	api.OKResponse(w, response)
}

func (s *CatalogService) GetAllCategories() ([]models.Category, error) {
	return s.CategoryRepo.GetAllCategories()
}

func (s *CatalogService) CreateCategory(code, name string) (*models.Category, error) {
	return s.CategoryRepo.CreateCategory(code, name)
}

func (s *CatalogService) GetAllProducts(params models.ProductQueryParams) ([]models.Product, error) {
	return s.ProductsRepo.GetAllProducts(params)
}

func (s *CatalogService) GetProductByCode(code string) (*models.Product, error) {
	return s.ProductsRepo.GetProductByCode(code)
}

// HELPERS
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
	if errors.Is(err, models.ErrInvalidQueryParam) {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, models.ErrCreationCategory) {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, models.ErrNotImplemented) {
		api.ErrorResponse(w, http.StatusNotImplemented, err.Error())
		return
	}

	// For all other errors
	api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
}

func parseQueryParams(r *http.Request) models.ProductQueryParams {
	queryParams := models.ProductQueryParams{}

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	categoryName := r.URL.Query().Get("category")
	priceParam := r.URL.Query().Get("price")

	queryParams.Limit = DefaultPageSize

	// Parse and validate limit
	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err == nil {
			if limit < MinPageSize {
				queryParams.Limit = MinPageSize
			} else if limit > MaxPageSize {
				queryParams.Limit = MaxPageSize
			} else {
				queryParams.Limit = limit
			}
		}
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offsetParam == "" || offset < 0 {
		queryParams.Offset = DefaultOffset
	} else {
		queryParams.Offset = max(offset, 0)
	}

	queryParams.PriceMax = maxPrice
	if priceParam != "" {
		queryParams.PriceMax, _ = strconv.ParseFloat(priceParam, 64)
	}

	queryParams.CategoryName = categoryName

	return queryParams
}
