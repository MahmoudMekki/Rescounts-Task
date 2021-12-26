package handler

import (
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
)

type MiddleWare struct {
	l            *log.Logger
	tokenService token.Service
}

func NewMiddleWare(l *log.Logger, tkn token.Service) *MiddleWare {
	return &MiddleWare{
		l:            l,
		tokenService: tkn,
	}
}

func (m *MiddleWare) MW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		token, ok := m.tokenService.ExtractToken(req)
		if !ok {
			rhttp.RespondJSON(rw, http.StatusUnauthorized, "unauthorized user")
			return
		} else {
			if !token.Verify() {
				rhttp.RespondJSON(rw, http.StatusUnauthorized, "unauthorized user")
				return
			}
		}
		next.ServeHTTP(rw, req)
	})
}
