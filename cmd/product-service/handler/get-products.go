package handler

import (
	"github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
)

type GetProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
}

func NewGetProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service) *GetProductHandler {
	return &GetProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
	}
}

func (product *GetProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	prods, err := product.productRepo.GetProducts()
	if err != nil {
		resp := data.GetProductsResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	resp := data.GetProductsResponse{Status: 1, Data: &data.GetProductsData{Products: prods}}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}
