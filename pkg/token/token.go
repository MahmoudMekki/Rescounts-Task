package token

import (
	"net/http"
	"time"
)

// Token define token info
type Token interface {
	Signed() string
	ExpireAt() int64
	UserID() int64
	IsAdmin() bool
	Verify() bool
}

// Service --
type Service interface {
	New(userID int64, expireIn time.Duration) (Token, error)
	NewAdminToken(userID int64, expireIn time.Duration) (Token, error)
	Parse(token string) (Token, error)
	ExtractToken(req *http.Request) (Token, bool)
}
