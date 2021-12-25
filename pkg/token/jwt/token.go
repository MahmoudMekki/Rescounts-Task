package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"

	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims --
type Claims struct {
	*jwt.StandardClaims
	UserID  int64 `json:"user_id"`
	IsAdmin bool  `json:"is_admin"`
}

// Token --
type Token struct {
	*jwt.Token
	signedString string
}

// Signed --
func (t Token) Signed() string {
	return t.signedString
}

// ExpireAt --
func (t Token) ExpireAt() int64 {
	if t.Claims == nil {
		return -1
	}
	c, ok := t.Claims.(*Claims)
	if !ok {
		return -1
	}
	return c.ExpiresAt
}

// UserID --
func (t Token) UserID() int64 {
	if t.Claims == nil {
		return -1
	}
	c, ok := t.Claims.(*Claims)
	if !ok {
		return -1
	}

	return c.UserID
}

// IsAdmin --
func (t Token) IsAdmin() bool {
	if t.Claims == nil {
		return false
	}
	c, ok := t.Claims.(*Claims)
	if !ok {
		return false
	}
	return c.IsAdmin
}

// Encrypt --
func (t Token) Encrypt(key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(t.signedString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// tokenService --
type tokenService struct {
	secretKey string
}

// New --
func New(secretKey string) token.Service {
	return &tokenService{
		secretKey: secretKey,
	}
}

// Sign --
func (s *tokenService) Sign(t token.Token) (token.Token, error) {
	jwtToken, ok := t.(Token)
	if !ok {
		return nil, fmt.Errorf("invalid token type, only support jwt.Token")
	}
	claim, ok := jwtToken.Claims.(Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claim type, only support jwt.Token")
	}
	return s.sign(&claim)
}

// New --
func (s *tokenService) New(userID int64, expireIn time.Duration) (token.Token, error) {
	claims := newJwtClaim(userID, false, expireIn)
	t, err := s.sign(claims)
	return t, err
}

// NewAdminToken --
func (s *tokenService) NewAdminToken(userID int64, expireIn time.Duration) (token.Token, error) {
	claims := newJwtClaim(userID, true, expireIn)
	return s.sign(claims)
}

// Parse --
func (s *tokenService) Parse(token string) (token.Token, error) {
	c := &Claims{}
	result := &Token{
		Token: &jwt.Token{
			Claims: c,
		},
	}
	_, err := jwt.ParseWithClaims(token,
		c, func(token *jwt.Token) (interface{}, error) {
			if c.UserID <= 0 {
				return nil, fmt.Errorf("invalid token, userID after parse:%v", c.UserID)
			}
			// token.Claims = c
			result.Token = token
			return []byte(s.secretKey), nil
		})
	result.signedString = token

	return result, err
}

func newJwtClaim(
	userID int64,
	isAdmin bool,
	expireIn time.Duration,
) *Claims {

	tokenExpireInUnix := time.Now().Add(expireIn).Unix()
	return &Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: tokenExpireInUnix,
		},
		UserID:  userID,
		IsAdmin: isAdmin,
	}
}

func (s *tokenService) sign(claims *Claims) (*Token, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := t.SignedString([]byte(s.secretKey))
	return &Token{
		Token:        t,
		signedString: signedString,
	}, err
}
