package token

import (
	"time"
)

// Token define token info
type Token interface {
	Signed() string
	ExpireAt() int64
	UserID() int64
	IsAdmin() bool
	Encrypt(key []byte) string
}

// Service --
type Service interface {
	New(userID int64, expireIn time.Duration) (Token, error)
	NewAdminToken(userID int64, expireIn time.Duration) (Token, error)
	Parse(token string) (Token, error)
}
