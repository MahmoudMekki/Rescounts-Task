package data

import "github.com/MahmoudMekki/Rescounts-Task/pkg/model"

// GetProductsResponse --
type GetProductsResponse struct {
	Status  int64            `json:"status"`
	Data    *GetProductsData `json:"data,omitempty"`
	Message string           `json:"message,omitempty"`
}

type GetProductsData struct {
	Products []model.Product `json:"products"`
}
