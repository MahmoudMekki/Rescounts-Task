package handler

import (
	"github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type DeleteProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
}

func NewDeleteProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service) *DeleteProductHandler {
	return &DeleteProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
	}
}

func (product *DeleteProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	token, _ := product.tokenService.ExtractToken(req)
	if !token.IsAdmin() {
		resp := data.DeleteProductResponse{Status: 0, Message: "Not allowed"}
		rhttp.RespondJSON(rw, http.StatusUnauthorized, resp)
		return
	}
	vars := mux.Vars(req)
	sProdID := vars["prod_id"]
	prodId, err := strconv.Atoi(sProdID)
	if err != nil {
		product.l.Println(err.Error())
		resp := data.DeleteProductResponse{Status: 0, Message: "bad product id"}
		rhttp.RespondJSON(rw, http.StatusBadRequest, resp)
		return
	}
	err = product.productRepo.DeleteProductByID(int64(prodId))
	if err != nil {
		resp := data.DeleteProductResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusBadRequest, resp)
		return
	}
	resp := data.DeleteProductResponse{Status: 1}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}
