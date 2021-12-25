package data

// UpdateProductRequest --
type UpdateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// UpdateProductResponse --
type UpdateProductResponse struct {
	Message   string `json:"message"`
	ProductID int64  `json:"product_id"`
}
