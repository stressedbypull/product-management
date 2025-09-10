package catalog

import (
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo models.ProductInterface
}

func NewCatalogHandler(r models.ProductInterface) *CatalogHandler {
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
		}
	}

	// Return the products as a JSON response
	response := Response{
		Products: products,
	}

	api.OKResponse(w, response)
}

// handleError checks for specific error types and returns appropriate HTTP status codes
func (h *CatalogHandler) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrProductNotFound) {
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
