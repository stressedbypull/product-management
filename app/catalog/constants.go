package catalog

// Error messages
const (
	ErrMsgInvalidRequest      = "invalid request"
	ErrMsgProductNotFound     = "product not found"
	ErrMsgCategoryNotFound    = "category not found"
	ErrMsgInternalServerError = "internal server error"
)

// Repository method names
const (
	GetAllProducts   = "GetAllProducts"
	GetAllCategories = "GetAllCategories"
)

// API Endpoints
const (
	EndpointCatalog    = "/catalog"
	EndpointCategories = "/categories"
)

// HTTP Methods
const (
	HTTPMethodGet  = "GET"
	HTTPMethodPost = "POST"
)

// Query Parameters
const (
	QueryParamPage     = "page"
	QueryParamPageSize = "page_size"
	QueryParamCategory = "category"
	QueryParamPriceMax = "price_max"
)

// Common Status Codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// Filter Options
const (
	FilterByCategory = "category"
	FilterByPriceMax = "price_max"
)

// Pagination Defaults
const (
	DefaultPageSize = 10  // Default limit if not specified
	MinPageSize     = 1   // Minimum allowed limit
	MaxPageSize     = 100 // Maximum allowed limit
	DefaultOffset   = 0   // Default offset if not specified
)
const (
	maxPrice = 10000.0
)
