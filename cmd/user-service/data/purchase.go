package data

// PurchaseResponse --
type PurchaseResponse struct {
	Status  int64         `json:"status"`
	Data    *PurchaseData `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
}
type PurchaseData struct {
	TransactionID string `json:"transaction_id"`
}
