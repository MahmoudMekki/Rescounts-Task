package data

// AddCardRequest --
type AddCardRequest struct {
	CardNumber  string `json:"card_number"`
	ExpireMonth string `json:"expire_month"`
	ExpireYear  string `json:"expire_year"`
	CVC         string `json:"cvc"`
}

// AddCardResponse --
type AddCardResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message,omitempty"`
}
