package data

// AddProductRequest --
type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// AddProductResponse --
type AddProductResponse struct {
	Status  int64           `json:"status"`
	Data    *AddProductData `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

type AddProductData struct {
	ProductID int64 `json:"product_id"`
}
