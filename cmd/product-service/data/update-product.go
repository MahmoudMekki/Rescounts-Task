package data

// UpdateProductRequest --
type UpdateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// UpdateProductResponse --
type UpdateProductResponse struct {
	Status  int64              `json:"status"`
	Data    *UpdateProductData `json:"data,omitempty"`
	Message string             `json:"message,omitempty"`
}
type UpdateProductData struct {
	ProductID int64 `json:"product_id"`
}
