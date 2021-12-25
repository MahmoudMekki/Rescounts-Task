package data

// AddProductRequest --
type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// AddProductResponse --
type AddProductResponse struct {
	Message   string `json:"message"`
	ProductID int64  `json:"product_id"`
}
