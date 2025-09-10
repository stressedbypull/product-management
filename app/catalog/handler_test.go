package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetProducts(t *testing.T) {
	//Available Handlers
	var (
		GetAllProducts = "GetAllProducts"
		//CreateProduct = "CreateProduct"
	)

	var expectedResp Response
	repoProducts := []models.Product{
		{Code: "P001", Price: decimal.NewFromFloat(10.0)},
		{Code: "P002", Price: decimal.NewFromFloat(20.0)},
	}

	expectedResp.Products = []Product{
		{Code: "P001", Price: 10.0},
		{Code: "P002", Price: 20.0},
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	t.Run("successful retrieval of products", func(t *testing.T) {
		mockProduct := new(mocks.MockProductRepo)
		mockProduct.On(GetAllProducts).Return(repoProducts, nil)

		handler := NewCatalogHandler(mockProduct)

		//make request
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		recorder := httptest.NewRecorder()
		handler.HandleRetrieveProducts(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp Response
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("product not found error", func(t *testing.T) {
		notFoundErr := models.ErrProductNotFound

		mockProduct := new(mocks.MockProductRepo)
		mockProduct.On(GetAllProducts).Return(nil, notFoundErr)

		handler := NewCatalogHandler(mockProduct)

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		recorder := httptest.NewRecorder()
		handler.HandleRetrieveProducts(recorder, req)

		// Assert results - should be 404 Not Found for this specific error
		resp := recorder.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Check response body
		var actualResp ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Contains(t, actualResp.Error, "product not found")

		// Check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

}

// -- IGNORE --
// setupMockRepo is a helper function to set up the mock repository
// repoInterface should be your interface type (e.g., models.ProductInterface)
// Usage example in a test
// mockRepo := new(MockProductRepo)
// setupMockRepo(*mockRepo, "GetAllProducts", sampleProducts, nil)
// ...call handler and assert results...
// func setupMockRepo(repo *mock.Mock, method string, returnValue any, err error) *mock.Mock {
// 	repo.On(method).Return(returnValue, err)
// 	return repo
// }
