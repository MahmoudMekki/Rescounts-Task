package data

// LoginRequest --
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse --
type LoginResponse struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
