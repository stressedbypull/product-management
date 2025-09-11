package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetProductsByID(t *testing.T) {
	//Available Handlers
	var (
		GetProductByCode = "GetProductByCode"
	)
	repoProduct := &models.Product{
		Code:     "P001",
		Price:    decimal.NewFromFloat(10.0),
		Category: models.Category{ID: 1, Code: "C001", Name: "Category 1"},
		Variants: []models.Variant{
			{ID: 1, Name: "Small", SKU: "P001-S", Price: decimal.NewFromFloat(8.0)},
			{ID: 2, Name: "Medium", SKU: "P001-M", Price: decimal.NewFromFloat(10.0)},
			{ID: 3, Name: "Large", SKU: "P001-L", Price: decimal.NewFromFloat(12.0)},
		},
	}
	expectedResp := Product{
		Code:  "P001",
		Price: 10.0,
		Category: Category{
			ID:   1,
			Code: "C001",
			Name: "Category 1",
		},
		Variants: []Variant{
			{ID: 1, Name: "Small", SKU: "P001-S", Price: 8.0},
			{ID: 2, Name: "Medium", SKU: "P001-M", Price: 10.0},
			{ID: 3, Name: "Large", SKU: "P001-L", Price: 12.0},
		},
	}

	t.Run("successful retrieval of product by code and its variants", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)
		mockProduct.On(GetProductByCode, "P001").Return(repoProduct, nil)

		handler := NewCatalogHandler(mockProduct)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /catalog/{code}", handler.HandleRetrieveProductByCode)

		//make request
		req := httptest.NewRequest(http.MethodGet, "/catalog/P001", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp Product
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})
}

func TestHandlerGetProducts(t *testing.T) {
	//Available Handlers
	var (
		GetAllProducts = "GetAllProducts"
	)

	repoProducts := []models.Product{
		{Code: "P001", Price: decimal.NewFromFloat(10.0), Category: models.Category{ID: 1, Code: "C001", Name: "Category 1"}},
		{Code: "P002", Price: decimal.NewFromFloat(20.0), Category: models.Category{ID: 2, Code: "C002", Name: "Category 2"}},
	}

	expectedResp := ProductsResponse{
		Total: len(repoProducts),
		Products: []Product{
			{Code: "P001", Price: 10.0, Category: Category{ID: 1, Code: "C001", Name: "Category 1"}},
			{Code: "P002", Price: 20.0, Category: Category{ID: 2, Code: "C002", Name: "Category 2"}},
		},
	}

	defaultParams := models.ProductQueryParams{
		Limit:        10,
		Offset:       0,
		PriceMax:     10000, // default max price
		CategoryName: "",
	}

	t.Run("successful retrieval of products, with limit and offset", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)
		params := defaultParams
		params.Limit = 1
		params.Offset = 1
		//call to the db
		mockProduct.On(GetAllProducts, params).Return(repoProducts, nil)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		//make request
		endpoint := EndpointCatalog + "?limit=" + strconv.Itoa(params.Limit) + "&offset=" + strconv.Itoa(params.Offset)
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp ProductsResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		assert.Equal(t, len(expectedResp.Products), actualResp.Total)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("successful retrieval limit and offset not specified, success", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)

		//call to the db
		mockProduct.On(GetAllProducts, defaultParams).Return(repoProducts, nil)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		//make request
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp ProductsResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("successful retrieval after filtering by price max", func(t *testing.T) {
		params := defaultParams
		params.PriceMax = 15.0
		mockProduct := new(mocks.MockCatalogRepository)
		filteredProducts := []models.Product{
			{Code: "P001", Price: decimal.NewFromFloat(10.0), Category: models.Category{ID: 1, Code: "C001", Name: "Category 1"}},
		}
		expectedResp := ProductsResponse{
			Total: len(filteredProducts),
			Products: []Product{
				{Code: "P001", Price: 10.0, Category: Category{ID: 1, Code: "C001", Name: "Category 1"}},
			},
		}

		//call to the db
		mockProduct.On(GetAllProducts, params).Return(filteredProducts, nil)

		handler := NewCatalogHandler(mockProduct)

		//make request
		req := httptest.NewRequest(http.MethodGet, "/catalog?price=15", nil)
		recorder := httptest.NewRecorder()
		handler.HandleRetrieveProducts(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp ProductsResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		assert.Equal(t, len(expectedResp.Products), actualResp.Total)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("Successful retrieval with category filter", func(t *testing.T) {
		params := defaultParams
		params.CategoryName = "Category 1"
		mockProduct := new(mocks.MockCatalogRepository)
		filteredProducts := []models.Product{
			{Code: "P001", Price: decimal.NewFromFloat(10.0), Category: models.Category{ID: 1, Code: "C001", Name: "Category 1"}},
		}
		expectedResp := ProductsResponse{
			Total: len(filteredProducts),
			Products: []Product{
				{Code: "P001", Price: 10.0, Category: Category{ID: 1, Code: "C001", Name: "Category 1"}},
			},
		}

		//call to the db
		mockProduct.On(GetAllProducts, params).Return(filteredProducts, nil)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		//make request
		baseURL := "/catalog"
		queryParams := url.Values{}
		queryParams.Add("category", "Category 1")
		fullURL := baseURL + "?" + queryParams.Encode()
		req := httptest.NewRequest(http.MethodGet, fullURL, nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp ProductsResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		assert.Equal(t, len(expectedResp.Products), actualResp.Total)
		assert.Equal(t, expectedResp.Products[0].Category.Name, actualResp.Products[0].Category.Name)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("product not found error", func(t *testing.T) {
		notFoundErr := models.ErrProductNotFound

		mockProduct := new(mocks.MockCatalogRepository)
		mockProduct.On(GetAllProducts, defaultParams).Return(nil, notFoundErr)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		// Assert results - should be 404 Not Found for this specific error
		resp := recorder.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var actualResp ErrorResponse
		// Check response body
		err := json.NewDecoder(resp.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Contains(t, actualResp.Error, "product not found")

		// Check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

}

func TestHandlerGetCategories(t *testing.T) {
	//Available Handlers
	var (
		GetAllCategories = "GetAllCategories"
		CreateCategory   = "CreateCategory"
	)

	repoCategories := []models.Category{
		{ID: 1, Code: "C001", Name: "Category 1"},
		{ID: 2, Code: "C002", Name: "Category 2"},
	}

	// 2. Use a dedicated response struct
	expectedResp := CategoriesResponse{
		Total: len(repoCategories),
		Categories: []Category{
			{ID: 1, Code: "C001", Name: "Category 1"},
			{ID: 2, Code: "C002", Name: "Category 2"},
		},
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	t.Run("successful retrieval of Categories", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)
		//call to the db
		mockProduct.On(GetAllCategories).Return(repoCategories, nil)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		//make request
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		//assert results, 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//check response body is as expected
		var actualResp CategoriesResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp, actualResp)
		//check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("Create category not implemented", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)

		// Setup mock expectations
		mockProduct.On(CreateCategory, "C003", "Category 3").Return(nil, models.ErrNotImplemented)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		reqBody := strings.NewReader(`{"code":"C003","name":"Category 3"}`)

		// Make request
		req := httptest.NewRequest(http.MethodPost, "/categories", reqBody)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		// Assert results - should be 501 Not Implemented
		resp := recorder.Result()
		assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)

		var actualResp ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Contains(t, actualResp.Error, "not implemented")

		// Check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("Create category success", func(t *testing.T) {
		mockProduct := new(mocks.MockCatalogRepository)

		newCategory := &models.Category{ID: 3, Code: "C003", Name: "Category 3"}
		// Setup mock expectations
		mockProduct.On(CreateCategory, "C003", "Category 3").Return(newCategory, nil)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		reqBody := strings.NewReader(`{"code":"C003","name":"Category 3"}`)

		// Make request
		req := httptest.NewRequest(http.MethodPost, "/categories", reqBody)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		// Assert results - should be 200 OK
		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actualResp struct {
			Category *models.Category `json:"category"`
		}
		err := json.NewDecoder(resp.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Equal(t, newCategory, actualResp.Category)

		// Check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

	t.Run("category not found error", func(t *testing.T) {
		notFoundErr := models.ErrCategoryNotFound

		mockProduct := new(mocks.MockCatalogRepository)
		mockProduct.On(GetAllCategories).Return(nil, notFoundErr)

		handler := NewCatalogHandler(mockProduct)
		mux := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		// Assert results - should be 404 Not Found for this specific error
		resp := recorder.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Check response body
		var actualResp ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&actualResp)
		assert.NoError(t, err)
		assert.Contains(t, actualResp.Error, "category not found")

		// Check that the mock was called
		t.Cleanup(func() {
			mockProduct.AssertExpectations(t)
		})
	})

}

// HELPERS
// Helper function to setup consistent router for all tests
func setupTestRouter(handler *CatalogHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", handler.HandleRetrieveProducts)
	mux.HandleFunc("GET /catalog/{code}", handler.HandleRetrieveProductByCode)
	mux.HandleFunc("GET /categories", handler.HandleRetrieveCategories)
	mux.HandleFunc("POST /categories", handler.HandleCreateCategory)
	return mux
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
