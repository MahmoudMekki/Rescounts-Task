package data

// CreateAccountRequest --
type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// CreateAccountResponse --
type CreateAccountResponse struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
